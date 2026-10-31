package mail

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"brightloop/internal/testdb"
)

func TestMain(m *testing.M) {
	if err := testdb.Setup("bl_test_mail"); err != nil {
		fmt.Println("testdb:", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func enqueue(t *testing.T, to string) {
	t.Helper()
	if err := g.DB().Transaction(context.Background(), func(ctx context.Context, tx gdb.TX) error {
		return Enqueue(ctx, tx, "", Email{To: to, Subject: "Hi", Body: "Hello"})
	}); err != nil {
		t.Fatal(err)
	}
}

func job(to string) gdb.Record {
	r, _ := g.DB().GetOne(context.Background(), `SELECT state, attempts, last_error,
		TIMESTAMPDIFF(SECOND, NOW(3), effective_at) AS wait_s FROM job_queue WHERE JSON_EXTRACT(payload, '$.to') = ?`, to)
	return r
}

func TestWorkSendsAndMarksDone(t *testing.T) {
	var sent []Email
	Send = func(_ context.Context, e Email) error { sent = append(sent, e); return nil }
	enqueue(t, "a@ok.test")
	enqueue(t, "b@ok.test")
	Work(context.Background())
	if len(sent) != 2 {
		t.Fatalf("sent %d emails, want 2", len(sent))
	}
	for _, to := range []string{"a@ok.test", "b@ok.test"} {
		if s := job(to)["state"].String(); s != "done" {
			t.Errorf("%s: state %s", to, s)
		}
	}
	Work(context.Background())
	if len(sent) != 2 {
		t.Fatal("a done email was sent again")
	}
}

func TestFailingSendRetriesWithBackoffThenFails(t *testing.T) {
	ctx := context.Background()
	Send = func(context.Context, Email) error { return errors.New("connection refused") }
	enqueue(t, "down@fail.test")
	g.DB().Exec(ctx, "UPDATE job_queue SET max_attempts = 2 WHERE JSON_EXTRACT(payload, '$.to') = 'down@fail.test'")

	Work(ctx)
	r := job("down@fail.test")
	if r["state"].String() != "available" || r["attempts"].Int() != 1 || r["last_error"].String() != "connection refused" {
		t.Fatalf("after 1st failure: %v", r)
	}
	if w := r["wait_s"].Int(); w < 110 || w > 120 { // 2^1 minutes
		t.Errorf("retry in %d s, want ~120", w)
	}
	Work(ctx) // not due yet: no attempt
	if job("down@fail.test")["attempts"].Int() != 1 {
		t.Fatal("retried before backoff")
	}
	g.DB().Exec(ctx, "UPDATE job_queue SET effective_at = NOW(3) WHERE JSON_EXTRACT(payload, '$.to') = 'down@fail.test'")
	Work(ctx)
	if r = job("down@fail.test"); r["state"].String() != "failed" || r["attempts"].Int() != 2 {
		t.Fatalf("after max attempts: %v", r)
	}
}

func TestRolledBackChangeQueuesNothing(t *testing.T) {
	boom := errors.New("business change failed")
	err := g.DB().Transaction(context.Background(), func(ctx context.Context, tx gdb.TX) error {
		if err := Enqueue(ctx, tx, "", Email{To: "never@rollback.test"}); err != nil {
			return err
		}
		return boom
	})
	if !errors.Is(err, boom) {
		t.Fatal(err)
	}
	if !job("never@rollback.test").IsEmpty() {
		t.Fatal("job survived the rollback")
	}
}

func TestSMTPRejectsHeaderInjection(t *testing.T) {
	if err := sendSMTP(context.Background(), Email{To: "a@b.test\r\nBcc: victim@x.test", Subject: "x"}); err == nil {
		t.Fatal("newline in To was accepted")
	}
}
