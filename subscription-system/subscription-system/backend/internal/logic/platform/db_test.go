package platform

import (
	"context"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func TestTxRetriesOnceOnDeadlock(t *testing.T) {
	for _, c := range []struct {
		failures, wantCalls int
		wantErr             bool
	}{{0, 1, false}, {1, 2, false}, {2, 2, true}} {
		calls := 0
		err := Tx(context.Background(), func(ctx context.Context, tx gdb.TX) error {
			calls++
			if calls <= c.failures {
				return &mysql.MySQLError{Number: 1213}
			}
			return nil
		})
		if calls != c.wantCalls || (err != nil) != c.wantErr {
			t.Errorf("%d failures: %d calls, err %v", c.failures, calls, err)
		}
	}
}

func TestAuditRecordsActor(t *testing.T) {
	ctx := context.Background()
	r, err := g.DB().Exec(ctx, "INSERT INTO users (email, password_hash, role) VALUES (?, 'x', 'admin')", unique("actor")+"@t.test")
	if err != nil {
		t.Fatal(err)
	}
	uid, _ := r.LastInsertId()
	actx := WithActor(ctx, Actor{UserID: uint64(uid), Role: "admin"})
	if err = Tx(actx, func(ctx context.Context, tx gdb.TX) error {
		return Audit(tx, "users", uint64(uid), "test-actor", "", nil)
	}); err != nil {
		t.Fatal(err)
	}
	if err = Tx(ctx, func(ctx context.Context, tx gdb.TX) error {
		return Audit(tx, "users", uint64(uid), "test-system", "", nil)
	}); err != nil {
		t.Fatal(err)
	}
	rows, _ := g.DB().GetAll(ctx, "SELECT action, changed_by FROM audit_log WHERE record_id = ? AND action LIKE 'test-%' ORDER BY id", uid)
	if len(rows) != 2 || rows[0]["changed_by"].Int64() != uid || !rows[1]["changed_by"].IsNil() {
		t.Fatalf("changed_by wrong: %v", rows)
	}
}
