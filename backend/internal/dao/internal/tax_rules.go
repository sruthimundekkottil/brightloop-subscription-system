// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TaxRulesDao is the data access object for table tax_rules.
type TaxRulesDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns TaxRulesColumns // columns contains all the column names of Table for convenient usage.
}

// TaxRulesColumns defines and stores column names for table tax_rules.
type TaxRulesColumns struct {
	Id            string //
	Name          string //
	RatePercent   string //
	IsDefault     string //
	IsActive      string //
	DefaultMarker string //
}

// taxRulesColumns holds the columns for table tax_rules.
var taxRulesColumns = TaxRulesColumns{
	Id:            "id",
	Name:          "name",
	RatePercent:   "rate_percent",
	IsDefault:     "is_default",
	IsActive:      "is_active",
	DefaultMarker: "default_marker",
}

// NewTaxRulesDao creates and returns a new DAO object for table data access.
func NewTaxRulesDao() *TaxRulesDao {
	return &TaxRulesDao{
		group:   "default",
		table:   "tax_rules",
		columns: taxRulesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TaxRulesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TaxRulesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TaxRulesDao) Columns() TaxRulesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TaxRulesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TaxRulesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TaxRulesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
