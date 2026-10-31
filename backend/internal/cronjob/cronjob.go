// Package cronjob registers scheduled jobs, as UniBee does in internal/cronjob.
package cronjob

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"

	"brightloop/internal/logic/billing"
	"brightloop/internal/logic/platform"
)

const BillingJob = "billing-nightly"

// Register adds the nightly billing run at 00:05 server time (UTC in our containers).
// AddSingleton skips a tick while the previous run is still going; billing.Run is safe to re-run regardless.
func Register(ctx context.Context) error {
	_, err := gcron.AddSingleton(ctx, "0 5 0 * * *", runBilling, BillingJob)
	return err
}

func runBilling(ctx context.Context) {
	today, err := platform.Today(ctx)
	if err != nil {
		g.Log().Error(ctx, "billing: reading today:", err)
		return
	}
	res, err := billing.Run(ctx, today, false)
	if err != nil {
		g.Log().Error(ctx, "billing run failed:", err)
	}
	if res != nil {
		g.Log().Infof(ctx, "billing run %d for %s: due %d, created %d, skipped %d, failed %d",
			res.RunID, today.Format(platform.DateLayout), res.Due, res.Created, res.Skipped, res.Failed)
	}
}
