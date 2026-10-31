package seed

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/crypto/bcrypt"

	"brightloop/internal/testdb"
)

func TestMain(m *testing.M) {
	if err := testdb.Setup("bl_test_seed"); err != nil {
		fmt.Println("testdb:", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestSeedIsRepeatable(t *testing.T) {
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		if err := Run(ctx); err != nil {
			t.Fatalf("run %d: %v", i+1, err)
		}
	}
	counts := map[string]int{"users": 3, "customers": 1, "staff_profiles": 1, "company_settings": 1, "payment_terms": 3, "close_reasons": 8}
	for table, want := range counts {
		if n, _ := g.DB().GetCount(ctx, "SELECT COUNT(*) FROM "+table); n != want {
			t.Errorf("%s: %d rows, want %d", table, n, want)
		}
	}
	hash, _ := g.DB().GetValue(ctx, "SELECT password_hash FROM users WHERE email = ?", AdminEmail)
	if bcrypt.CompareHashAndPassword(hash.Bytes(), []byte(DevPassword)) != nil {
		t.Error("admin password hash does not match DevPassword")
	}
	if n, _ := g.DB().GetCount(ctx, "SELECT COUNT(*) FROM close_reasons WHERE name = 'Renewed' AND is_churn = FALSE"); n != 1 {
		t.Error("Renewed must not count as churn")
	}
}
