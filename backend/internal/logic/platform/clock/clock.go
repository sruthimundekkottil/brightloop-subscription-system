// Package clock is the team-facing name for the platform clock. Every date decision goes through Today.
package clock

import (
	"context"
	"time"

	"brightloop/internal/logic/platform"
)

// Today is the UTC business date; in development it honours the forward-only test clock.
func Today(ctx context.Context) (time.Time, error) { return platform.Today(ctx) }

// SetTestClock moves the development clock forward; see platform.SetTestClock.
func SetTestClock(ctx context.Context, date time.Time) error { return platform.SetTestClock(ctx, date) }
