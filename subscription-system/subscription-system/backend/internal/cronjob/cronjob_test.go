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
	defer gcron.Remove(MailJob)
	for _, name := range []string{BillingJob, MailJob} {
		e := gcron.Search(name)
		if e == nil {
			t.Fatal(name, "not registered")
		}
		if !e.IsSingleton() {
			t.Fatal(name, "must be a singleton")
		}
	}
	if err := Register(context.Background()); err == nil {
		t.Fatal("registering twice should fail on the duplicate name")
	}
}
