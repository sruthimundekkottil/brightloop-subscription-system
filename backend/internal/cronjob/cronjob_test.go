package cronjob

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/os/gcron"
)

func TestRegisterAddsNightlySingleton(t *testing.T) {
	if err := Register(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer gcron.Remove(BillingJob)
	e := gcron.Search(BillingJob)
	if e == nil {
		t.Fatal("billing job not registered")
	}
	if !e.IsSingleton() {
		t.Fatal("billing job must be a singleton")
	}
	if err := Register(context.Background()); err == nil {
		t.Fatal("registering twice should fail on the duplicate name")
	}
}
