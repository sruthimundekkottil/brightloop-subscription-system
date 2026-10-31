// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BillingRunsDao is the data access object for table billing_runs.
type BillingRunsDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns BillingRunsColumns // columns contains all the column names of Table for convenient usage.
}

// BillingRunsColumns defines and stores column names for table billing_runs.
type BillingRunsColumns struct {
	Id           string //
	AsOfDate     string //
	TriggeredBy  string //
	DryRun       string //
	StartedAt    string //
	FinishedAt   string //
	DueCount     string //
	CreatedCount string //
	SkippedCount string //
	FailedCount  string //
	Errors       string //
}

// billingRunsColumns holds the columns for table billing_runs.
var billingRunsColumns = BillingRunsColumns{
	Id:           "id",
	AsOfDate:     "as_of_date",
	TriggeredBy:  "triggered_by",
	DryRun:       "dry_run",
	StartedAt:    "started_at",
	FinishedAt:   "finished_at",
	DueCount:     "due_count",
	CreatedCount: "created_count",
	SkippedCount: "skipped_count",
	FailedCount:  "failed_count",
	Errors:       "errors",
}

// NewBillingRunsDao creates and returns a new DAO object for table data access.
func NewBillingRunsDao() *BillingRunsDao {
	return &BillingRunsDao{
		group:   "default",
		table:   "billing_runs",
		columns: billingRunsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *BillingRunsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *BillingRunsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *BillingRunsDao) Columns() BillingRunsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *BillingRunsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *BillingRunsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *BillingRunsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
