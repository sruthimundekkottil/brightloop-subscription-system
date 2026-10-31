package platform

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// JobLease is how long a claimed job stays reserved; after it, a crashed worker's job is claimable again.
const JobLease = 5 * time.Minute

var ErrLeaseLost = errors.New("job is no longer held by this worker")

type Job struct {
	ID       uint64
	Payload  string // JSON
	Attempts int    // including the current one
}

// Enqueue adds a job in the caller's transaction, so it exists only if the business change commits.
// A repeated non-empty dedupeKey is ignored.
func Enqueue(tx gdb.TX, queue, dedupeKey string, payload any, effectiveAt time.Time) error {
	js, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO job_queue (queue_name, dedupe_key, payload, effective_at)
		VALUES (?, NULLIF(?, ''), ?, ?) ON DUPLICATE KEY UPDATE id = id`,
		queue, dedupeKey, string(js), effectiveAt.UTC().Format("2006-01-02 15:04:05.000"))
	return err
}

// Claim reserves up to limit due jobs for owner. SKIP LOCKED lets concurrent workers take different rows without waiting.
func Claim(ctx context.Context, queue, owner string, limit int) ([]Job, error) {
	var jobs []Job
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		rows, err := tx.GetAll(`SELECT id, payload, attempts FROM job_queue
			WHERE queue_name = ? AND effective_at <= NOW(3)
			  AND (state = 'available' OR (state = 'processing' AND locked_until < NOW(3)))
			ORDER BY effective_at, id LIMIT ? FOR UPDATE SKIP LOCKED`, queue, limit)
		if err != nil || len(rows) == 0 {
			return err
		}
		ids := make([]uint64, len(rows))
		for i, r := range rows {
			ids[i] = r["id"].Uint64()
			jobs = append(jobs, Job{ID: ids[i], Payload: r["payload"].String(), Attempts: r["attempts"].Int() + 1})
		}
		_, err = tx.Exec(`UPDATE job_queue SET state = 'processing', locked_by = ?,
			locked_until = NOW(3) + INTERVAL ? SECOND, attempts = attempts + 1 WHERE id IN (?)`,
			owner, int(JobLease.Seconds()), ids)
		return err
	})
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// Done marks a claimed job finished.
func Done(ctx context.Context, id uint64, owner string) error {
	return heldUpdate(ctx, `UPDATE job_queue SET state = 'done', processed_at = NOW(3), locked_by = NULL, locked_until = NULL
		WHERE id = ? AND locked_by = ? AND state = 'processing'`, id, owner)
}

// Fail releases a claimed job for a retry after 2^attempts minutes (2, 4, 8, 16), or marks it failed once max_attempts is reached.
func Fail(ctx context.Context, id uint64, owner string, cause error) error {
	return heldUpdate(ctx, `UPDATE job_queue SET state = IF(attempts >= max_attempts, 'failed', 'available'),
		effective_at = NOW(3) + INTERVAL POW(2, attempts) MINUTE, locked_by = NULL, locked_until = NULL, last_error = ?
		WHERE id = ? AND locked_by = ? AND state = 'processing'`, cause.Error(), id, owner)
}

func heldUpdate(ctx context.Context, sql string, args ...any) error {
	res, err := g.DB().Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrLeaseLost
	}
	return nil
}
