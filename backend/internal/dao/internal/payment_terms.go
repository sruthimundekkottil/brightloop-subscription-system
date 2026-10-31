// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentTermsDao is the data access object for table payment_terms.
type PaymentTermsDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns PaymentTermsColumns // columns contains all the column names of Table for convenient usage.
}

// PaymentTermsColumns defines and stores column names for table payment_terms.
type PaymentTermsColumns struct {
	Id      string //
	Name    string //
	DueDays string //
}

// paymentTermsColumns holds the columns for table payment_terms.
var paymentTermsColumns = PaymentTermsColumns{
	Id:      "id",
	Name:    "name",
	DueDays: "due_days",
}

// NewPaymentTermsDao creates and returns a new DAO object for table data access.
func NewPaymentTermsDao() *PaymentTermsDao {
	return &PaymentTermsDao{
		group:   "default",
		table:   "payment_terms",
		columns: paymentTermsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PaymentTermsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PaymentTermsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PaymentTermsDao) Columns() PaymentTermsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PaymentTermsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PaymentTermsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PaymentTermsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
