// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductVariantValuesDao is the data access object for table product_variant_values.
type ProductVariantValuesDao struct {
	table   string                      // table is the underlying table name of the DAO.
	group   string                      // group is the database configuration group name of current DAO.
	columns ProductVariantValuesColumns // columns contains all the column names of Table for convenient usage.
}

// ProductVariantValuesColumns defines and stores column names for table product_variant_values.
type ProductVariantValuesColumns struct {
	VariantId               string //
	ProductAttributeValueId string //
}

// productVariantValuesColumns holds the columns for table product_variant_values.
var productVariantValuesColumns = ProductVariantValuesColumns{
	VariantId:               "variant_id",
	ProductAttributeValueId: "product_attribute_value_id",
}

// NewProductVariantValuesDao creates and returns a new DAO object for table data access.
func NewProductVariantValuesDao() *ProductVariantValuesDao {
	return &ProductVariantValuesDao{
		group:   "default",
		table:   "product_variant_values",
		columns: productVariantValuesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ProductVariantValuesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ProductVariantValuesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ProductVariantValuesDao) Columns() ProductVariantValuesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ProductVariantValuesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ProductVariantValuesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ProductVariantValuesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
