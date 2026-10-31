package platform

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"brightloop/internal/testdb"
)

func TestMain(m *testing.M) {
	if err := testdb.Setup("bl_test_platform"); err != nil {
		fmt.Println("testdb:", err)
		os.Exit(1)
	}
	if _, err := g.DB().Exec(context.Background(), "INSERT INTO company_settings (id) VALUES (1)"); err != nil {
		fmt.Println("seed:", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func tx(t *testing.T, f func(tx gdb.TX) error) error {
	t.Helper()
	return g.DB().Transaction(context.Background(), func(ctx context.Context, tx gdb.TX) error { return f(tx) })
}

func TestTestClock(t *testing.T) {
	ctx := context.Background()
	realToday := time.Now().UTC().Truncate(24 * time.Hour)
	future := realToday.AddDate(0, 2, 0)

	t.Setenv("APP_ENV", "production")
	if err := SetTestClock(ctx, future); !errors.Is(err, ErrClockDisabled) {
		t.Fatalf("outside development: got %v", err)
	}

	t.Setenv("APP_ENV", "development")
	if err := SetTestClock(ctx, future); err != nil {
		t.Fatal(err)
	}
	if got, _ := Today(ctx); !got.Equal(future) {
		t.Fatalf("Today = %s, want %s", got, future)
	}
	for _, d := range []time.Time{future, future.AddDate(0, 0, -1)} {
		if err := SetTestClock(ctx, d); !errors.Is(err, ErrClockBackwards) {
			t.Fatalf("set to %s: got %v, want ErrClockBackwards", d.Format(DateLayout), err)
		}
	}

	t.Setenv("APP_ENV", "production") // the stored clock is ignored outside development
	if got, _ := Today(ctx); !got.Equal(realToday) {
		t.Fatalf("production Today = %s, want %s", got, realToday)
	}
}

func TestNextNumberIsGaplessAcrossRollback(t *testing.T) {
	var first, second uint64
	_ = tx(t, func(tx gdb.TX) error { first, _ = NextNumber(tx, "invoice-test"); return nil })
	rollback := errors.New("rollback")
	if err := tx(t, func(tx gdb.TX) error { NextNumber(tx, "invoice-test"); return rollback }); err != rollback {
		t.Fatal(err)
	}
	_ = tx(t, func(tx gdb.TX) error { second, _ = NextNumber(tx, "invoice-test"); return nil })
	if first != 1 || second != 2 {
		t.Fatalf("got %d then %d, want 1 then 2", first, second)
	}
}

func TestAudit(t *testing.T) {
	if err := tx(t, func(tx gdb.TX) error {
		return Audit(tx, "invoices", 7, "create", "", map[string][2]any{"status": {nil, "draft"}})
	}); err != nil {
		t.Fatal(err)
	}
	r, err := g.DB().GetOne(context.Background(), "SELECT action, reason, changed_by, changes FROM audit_log WHERE table_name = 'invoices' AND record_id = 7")
	if err != nil {
		t.Fatal(err)
	}
	if r["action"].String() != "create" || !r["reason"].IsNil() || !r["changed_by"].IsNil() || r["changes"].String() != `{"status": [null, "draft"]}` {
		t.Fatalf("unexpected row %v", r.Map())
	}
}

func TestEnqueueDedupes(t *testing.T) {
	for i := 0; i < 2; i++ {
		if err := tx(t, func(tx gdb.TX) error {
			return Enqueue(tx, "dedupe", "invoice-sent:42", g.Map{"invoice_id": 42}, time.Now())
		}); err != nil {
			t.Fatal(err)
		}
	}
	if n, _ := g.DB().GetValue(context.Background(), "SELECT COUNT(*) FROM job_queue WHERE queue_name = 'dedupe'"); n.Int() != 1 {
		t.Fatalf("got %d jobs, want 1", n.Int())
	}
}

func TestTwoWorkersClaimDisjointJobs(t *testing.T) {
	ctx := context.Background()
	const total = 40
	_ = tx(t, func(tx gdb.TX) error {
		for i := 0; i < total; i++ {
			if err := Enqueue(tx, "race", "", g.Map{"n": i}, time.Now().Add(-time.Second)); err != nil {
				return err
			}
		}
		return nil
	})
	var mu sync.Mutex
	seen := map[uint64]string{}
	var wg sync.WaitGroup
	for _, w := range []string{"w1", "w2"} {
		wg.Add(1)
		go func(owner string) {
			defer wg.Done()
			for {
				jobs, err := Claim(ctx, "race", owner, 3)
				if err != nil {
					t.Error(err)
					return
				}
				if len(jobs) == 0 {
					return
				}
				for _, j := range jobs {
					mu.Lock()
					if prev, dup := seen[j.ID]; dup {
						t.Errorf("job %d claimed by %s and %s", j.ID, prev, owner)
					}
					seen[j.ID] = owner
					mu.Unlock()
					if err := Done(ctx, j.ID, owner); err != nil {
						t.Error(err)
					}
				}
			}
		}(w)
	}
	wg.Wait()
	if len(seen) != total {
		t.Fatalf("claimed %d jobs, want %d", len(seen), total)
	}
}

func TestFailRetriesThenGivesUp(t *testing.T) {
	ctx := context.Background()
	_ = tx(t, func(tx gdb.TX) error { return Enqueue(tx, "flaky", "", g.Map{}, time.Now().Add(-time.Second)) })
	state := func() (string, int) {
		r, _ := g.DB().GetOne(ctx, "SELECT state, attempts FROM job_queue WHERE queue_name = 'flaky'")
		return r["state"].String(), r["attempts"].Int()
	}
	makeDue := func() { g.DB().Exec(ctx, "UPDATE job_queue SET effective_at = NOW(3) - INTERVAL 1 SECOND WHERE queue_name = 'flaky'") }

	for attempt := 1; attempt <= 5; attempt++ { // max_attempts defaults to 5
		jobs, err := Claim(ctx, "flaky", "w", 1)
		if err != nil || len(jobs) != 1 || jobs[0].Attempts != attempt {
			t.Fatalf("attempt %d: jobs=%v err=%v", attempt, jobs, err)
		}
		if err := Done(ctx, jobs[0].ID, "someone-else"); !errors.Is(err, ErrLeaseLost) {
			t.Fatalf("Done by another owner: got %v", err)
		}
		if err := Fail(ctx, jobs[0].ID, "w", errors.New("smtp down")); err != nil {
			t.Fatal(err)
		}
		want := "available"
		if attempt == 5 {
			want = "failed"
		}
		if s, n := state(); s != want || n != attempt {
			t.Fatalf("after attempt %d: state=%s attempts=%d, want %s", attempt, s, n, want)
		}
		if jobs, _ := Claim(ctx, "flaky", "w", 1); len(jobs) != 0 {
			t.Fatalf("attempt %d: job claimable before its retry time", attempt)
		}
		makeDue()
	}
	if jobs, _ := Claim(ctx, "flaky", "w", 1); len(jobs) != 0 {
		t.Fatal("failed job was claimed again")
	}
}

func TestExpiredLeaseIsReclaimed(t *testing.T) {
	ctx := context.Background()
	_ = tx(t, func(tx gdb.TX) error { return Enqueue(tx, "crash", "", g.Map{}, time.Now().Add(-time.Second)) })
	jobs, _ := Claim(ctx, "crash", "dead-worker", 1)
	if len(jobs) != 1 {
		t.Fatal("no job claimed")
	}
	g.DB().Exec(ctx, "UPDATE job_queue SET locked_until = NOW(3) - INTERVAL 1 SECOND WHERE id = ?", jobs[0].ID)
	again, _ := Claim(ctx, "crash", "new-worker", 1)
	if len(again) != 1 || again[0].ID != jobs[0].ID || again[0].Attempts != 2 {
		t.Fatalf("reclaim got %v", again)
	}
	if err := Done(ctx, jobs[0].ID, "dead-worker"); !errors.Is(err, ErrLeaseLost) {
		t.Fatalf("stale owner Done: got %v", err)
	}
}
