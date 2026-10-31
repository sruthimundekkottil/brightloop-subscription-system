// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DiscountRulesDao is the data access object for table discount_rules.
type DiscountRulesDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns DiscountRulesColumns // columns contains all the column names of Table for convenient usage.
}

// DiscountRulesColumns defines and stores column names for table discount_rules.
type DiscountRulesColumns struct {
	Id                       string //
	Name                     string //
	DiscountType             string //
	Value                    string //
	Scope                    string //
	CustomerId               string //
	MinPurchaseAmount        string //
	MinQuantity              string //
	StartsOn                 string //
	EndsOn                   string //
	MaxUses                  string //
	MaxUsesPerCustomer       string //
	MaxCyclesPerSubscription string //
	MaxDiscountAmount        string //
	TimesUsed                string //
	IsActive                 string //
	CreatedBy                string //
	CreatedAt                string //
}

// discountRulesColumns holds the columns for table discount_rules.
var discountRulesColumns = DiscountRulesColumns{
	Id:                       "id",
	Name:                     "name",
	DiscountType:             "discount_type",
	Value:                    "value",
	Scope:                    "scope",
	CustomerId:               "customer_id",
	MinPurchaseAmount:        "min_purchase_amount",
	MinQuantity:              "min_quantity",
	StartsOn:                 "starts_on",
	EndsOn:                   "ends_on",
	MaxUses:                  "max_uses",
	MaxUsesPerCustomer:       "max_uses_per_customer",
	MaxCyclesPerSubscription: "max_cycles_per_subscription",
	MaxDiscountAmount:        "max_discount_amount",
	TimesUsed:                "times_used",
	IsActive:                 "is_active",
	CreatedBy:                "created_by",
	CreatedAt:                "created_at",
}

// NewDiscountRulesDao creates and returns a new DAO object for table data access.
func NewDiscountRulesDao() *DiscountRulesDao {
	return &DiscountRulesDao{
		group:   "default",
		table:   "discount_rules",
		columns: discountRulesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *DiscountRulesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *DiscountRulesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *DiscountRulesDao) Columns() DiscountRulesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *DiscountRulesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *DiscountRulesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *DiscountRulesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
