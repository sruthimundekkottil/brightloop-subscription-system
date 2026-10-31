// Package sequence hands out gapless document numbers inside the caller's transaction.
package sequence

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"

	"brightloop/internal/logic/platform"
)

// Next takes the next value of counter name. Counter names are "<document>-<year>", e.g. "invoice-2026",
// "subscription-2026"; billing.Run already uses "invoice-YYYY", so invoice confirmation must too.
func Next(ctx context.Context, tx gdb.TX, name string) (uint64, error) {
	return platform.NextNumber(tx.Ctx(ctx), name)
}

// Format renders a document number such as INV-2026-00001.
func Format(prefix string, year int, n uint64) string { return fmt.Sprintf("%s-%d-%05d", prefix, year, n) }
