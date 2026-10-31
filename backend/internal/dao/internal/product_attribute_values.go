// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductAttributeValuesDao is the data access object for table product_attribute_values.
type ProductAttributeValuesDao struct {
	table   string                        // table is the underlying table name of the DAO.
	group   string                        // group is the database configuration group name of current DAO.
	columns ProductAttributeValuesColumns // columns contains all the column names of Table for convenient usage.
}

// ProductAttributeValuesColumns defines and stores column names for table product_attribute_values.
type ProductAttributeValuesColumns struct {
	Id               string //
	ProductId        string //
	AttributeValueId string //
	PriceExtra       string //
	IsActive         string //
}

// productAttributeValuesColumns holds the columns for table product_attribute_values.
var productAttributeValuesColumns = ProductAttributeValuesColumns{
	Id:               "id",
	ProductId:        "product_id",
	AttributeValueId: "attribute_value_id",
	PriceExtra:       "price_extra",
	IsActive:         "is_active",
}

// NewProductAttributeValuesDao creates and returns a new DAO object for table data access.
func NewProductAttributeValuesDao() *ProductAttributeValuesDao {
	return &ProductAttributeValuesDao{
		group:   "default",
		table:   "product_attribute_values",
		columns: productAttributeValuesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ProductAttributeValuesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ProductAttributeValuesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ProductAttributeValuesDao) Columns() ProductAttributeValuesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ProductAttributeValuesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ProductAttributeValuesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ProductAttributeValuesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
