// Package audit records who changed what, in the same transaction as the change.
package audit

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"brightloop/internal/logic/platform"
)

// Log writes one audit_log row. changed_by is the logged-in user in ctx (NULL for system jobs);
// changes is usually {"field": [old, new]}.
func Log(ctx context.Context, tx gdb.TX, table string, id uint64, action, reason string, changes any) error {
	return platform.Audit(tx.Ctx(ctx), table, id, action, reason, changes)
}

// Update locks row id of table, writes only the columns in set whose value differs, and logs
// {"column": [old, new]} as action "update". A nil value means NULL. It reports false if the row does not exist.
// Column names come from code, never from the request. Pass money as StringFixed(2) so it compares with MySQL text.
func Update(ctx context.Context, tx gdb.TX, table string, id uint64, set map[string]any) (bool, error) {
	cols := make([]string, 0, len(set))
	for c := range set {
		cols = append(cols, c)
	}
	if len(cols) == 0 {
		n, err := tx.Ctx(ctx).GetCount(fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE id = ?", table), id)
		return n == 1, err
	}
	old, err := tx.Ctx(ctx).GetOne(fmt.Sprintf("SELECT `%s` FROM `%s` WHERE id = ? FOR UPDATE", strings.Join(cols, "`, `"), table), id)
	if err != nil || old.IsEmpty() {
		return false, err
	}
	changes := g.Map{}
	var assigns []string
	var args []any
	for _, c := range cols {
		v := set[c]
		if b, ok := v.(bool); ok { // MySQL returns BOOLEAN as 0/1
			v = map[bool]int{false: 0, true: 1}[b]
		}
		if old[c].IsNil() && v == nil || !old[c].IsNil() && v != nil && old[c].String() == fmt.Sprint(v) {
			continue
		}
		changes[c] = []any{old[c].Val(), v}
		assigns = append(assigns, "`"+c+"` = ?")
		args = append(args, v)
	}
	if len(assigns) == 0 {
		return true, nil
	}
	if _, err = tx.Ctx(ctx).Exec(fmt.Sprintf("UPDATE `%s` SET %s WHERE id = ?", table, strings.Join(assigns, ", ")), append(args, id)...); err != nil {
		return true, err
	}
	return true, Log(ctx, tx, table, id, "update", "", changes)
}
