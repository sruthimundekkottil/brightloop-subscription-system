// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DiscountApplicationsDao is the data access object for table discount_applications.
type DiscountApplicationsDao struct {
	table   string                      // table is the underlying table name of the DAO.
	group   string                      // group is the database configuration group name of current DAO.
	columns DiscountApplicationsColumns // columns contains all the column names of Table for convenient usage.
}

// DiscountApplicationsColumns defines and stores column names for table discount_applications.
type DiscountApplicationsColumns struct {
	Id             string //
	DiscountRuleId string //
	InvoiceId      string //
	CustomerId     string //
	SubscriptionId string //
	Amount         string //
	ReversedAt     string //
	AppliedAt      string //
}

// discountApplicationsColumns holds the columns for table discount_applications.
var discountApplicationsColumns = DiscountApplicationsColumns{
	Id:             "id",
	DiscountRuleId: "discount_rule_id",
	InvoiceId:      "invoice_id",
	CustomerId:     "customer_id",
	SubscriptionId: "subscription_id",
	Amount:         "amount",
	ReversedAt:     "reversed_at",
	AppliedAt:      "applied_at",
}

// NewDiscountApplicationsDao creates and returns a new DAO object for table data access.
func NewDiscountApplicationsDao() *DiscountApplicationsDao {
	return &DiscountApplicationsDao{
		group:   "default",
		table:   "discount_applications",
		columns: discountApplicationsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *DiscountApplicationsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *DiscountApplicationsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *DiscountApplicationsDao) Columns() DiscountApplicationsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *DiscountApplicationsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *DiscountApplicationsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *DiscountApplicationsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
