// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SubscriptionStatusHistoryDao is the data access object for table subscription_status_history.
type SubscriptionStatusHistoryDao struct {
	table   string                           // table is the underlying table name of the DAO.
	group   string                           // group is the database configuration group name of current DAO.
	columns SubscriptionStatusHistoryColumns // columns contains all the column names of Table for convenient usage.
}

// SubscriptionStatusHistoryColumns defines and stores column names for table subscription_status_history.
type SubscriptionStatusHistoryColumns struct {
	Id             string //
	SubscriptionId string //
	FromStatus     string //
	ToStatus       string //
	ChangedBy      string //
	Note           string //
	ChangedAt      string //
}

// subscriptionStatusHistoryColumns holds the columns for table subscription_status_history.
var subscriptionStatusHistoryColumns = SubscriptionStatusHistoryColumns{
	Id:             "id",
	SubscriptionId: "subscription_id",
	FromStatus:     "from_status",
	ToStatus:       "to_status",
	ChangedBy:      "changed_by",
	Note:           "note",
	ChangedAt:      "changed_at",
}

// NewSubscriptionStatusHistoryDao creates and returns a new DAO object for table data access.
func NewSubscriptionStatusHistoryDao() *SubscriptionStatusHistoryDao {
	return &SubscriptionStatusHistoryDao{
		group:   "default",
		table:   "subscription_status_history",
		columns: subscriptionStatusHistoryColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SubscriptionStatusHistoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SubscriptionStatusHistoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SubscriptionStatusHistoryDao) Columns() SubscriptionStatusHistoryColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SubscriptionStatusHistoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SubscriptionStatusHistoryDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SubscriptionStatusHistoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
