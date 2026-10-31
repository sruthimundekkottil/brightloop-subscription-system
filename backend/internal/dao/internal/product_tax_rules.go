// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductTaxRulesDao is the data access object for table product_tax_rules.
type ProductTaxRulesDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of current DAO.
	columns ProductTaxRulesColumns // columns contains all the column names of Table for convenient usage.
}

// ProductTaxRulesColumns defines and stores column names for table product_tax_rules.
type ProductTaxRulesColumns struct {
	ProductId string //
	TaxRuleId string //
}

// productTaxRulesColumns holds the columns for table product_tax_rules.
var productTaxRulesColumns = ProductTaxRulesColumns{
	ProductId: "product_id",
	TaxRuleId: "tax_rule_id",
}

// NewProductTaxRulesDao creates and returns a new DAO object for table data access.
func NewProductTaxRulesDao() *ProductTaxRulesDao {
	return &ProductTaxRulesDao{
		group:   "default",
		table:   "product_tax_rules",
		columns: productTaxRulesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ProductTaxRulesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ProductTaxRulesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ProductTaxRulesDao) Columns() ProductTaxRulesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ProductTaxRulesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ProductTaxRulesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ProductTaxRulesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
