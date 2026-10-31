// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BillingRecordsDao is the data access object for table billing_records.
type BillingRecordsDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns BillingRecordsColumns // columns contains all the column names of Table for convenient usage.
}

// BillingRecordsColumns defines and stores column names for table billing_records.
type BillingRecordsColumns struct {
	Id             string //
	SubscriptionId string //
	PeriodStart    string //
	PeriodEnd      string //
	InvoiceId      string //
	Inputs         string //
	CreatedAt      string //
}

// billingRecordsColumns holds the columns for table billing_records.
var billingRecordsColumns = BillingRecordsColumns{
	Id:             "id",
	SubscriptionId: "subscription_id",
	PeriodStart:    "period_start",
	PeriodEnd:      "period_end",
	InvoiceId:      "invoice_id",
	Inputs:         "inputs",
	CreatedAt:      "created_at",
}

// NewBillingRecordsDao creates and returns a new DAO object for table data access.
func NewBillingRecordsDao() *BillingRecordsDao {
	return &BillingRecordsDao{
		group:   "default",
		table:   "billing_records",
		columns: billingRecordsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *BillingRecordsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *BillingRecordsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *BillingRecordsDao) Columns() BillingRecordsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *BillingRecordsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *BillingRecordsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *BillingRecordsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
