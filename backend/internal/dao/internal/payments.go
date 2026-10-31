// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentsDao is the data access object for table payments.
type PaymentsDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns PaymentsColumns // columns contains all the column names of Table for convenient usage.
}

// PaymentsColumns defines and stores column names for table payments.
type PaymentsColumns struct {
	Id             string //
	InvoiceId      string //
	Method         string //
	Amount         string //
	PaymentDate    string //
	Note           string //
	IdempotencyKey string //
	RecordedBy     string //
	CreatedAt      string //
}

// paymentsColumns holds the columns for table payments.
var paymentsColumns = PaymentsColumns{
	Id:             "id",
	InvoiceId:      "invoice_id",
	Method:         "method",
	Amount:         "amount",
	PaymentDate:    "payment_date",
	Note:           "note",
	IdempotencyKey: "idempotency_key",
	RecordedBy:     "recorded_by",
	CreatedAt:      "created_at",
}

// NewPaymentsDao creates and returns a new DAO object for table data access.
func NewPaymentsDao() *PaymentsDao {
	return &PaymentsDao{
		group:   "default",
		table:   "payments",
		columns: paymentsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PaymentsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PaymentsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PaymentsDao) Columns() PaymentsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PaymentsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PaymentsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PaymentsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
