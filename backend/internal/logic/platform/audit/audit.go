// Package audit records who changed what, in the same transaction as the change.
package audit

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"

	"brightloop/internal/logic/platform"
)

// Log writes one audit_log row. changed_by is the logged-in user in ctx (NULL for system jobs);
// changes is usually {"field": [old, new]}.
func Log(ctx context.Context, tx gdb.TX, table string, id uint64, action, reason string, changes any) error {
	return platform.Audit(tx.Ctx(ctx), table, id, action, reason, changes)
}
