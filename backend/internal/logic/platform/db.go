package platform

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IsRetryable reports a deadlock (1213) or lock wait timeout (1205): the transaction was rolled back and may be re-run.
func IsRetryable(err error) bool {
	var me *mysql.MySQLError
	return errors.As(err, &me) && (me.Number == 1213 || me.Number == 1205)
}

// Tx runs f in a transaction and re-runs it once if MySQL rolled it back for a deadlock or lock wait timeout.
// f must not keep state between attempts.
func Tx(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	err := g.DB().Transaction(ctx, f)
	if IsRetryable(err) {
		err = g.DB().Transaction(ctx, f)
	}
	return err
}
