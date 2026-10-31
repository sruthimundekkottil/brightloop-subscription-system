// Package jobs enqueues background work (emails, reminders) in the caller's transaction.
package jobs

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gdb"

	"brightloop/internal/logic/platform"
)

// Enqueue adds a job that exists only if tx commits. A repeated non-empty dedupeKey is ignored.
func Enqueue(ctx context.Context, tx gdb.TX, queue, dedupeKey string, payload any, effectiveAt time.Time) error {
	return platform.Enqueue(tx.Ctx(ctx), queue, dedupeKey, payload, effectiveAt)
}
