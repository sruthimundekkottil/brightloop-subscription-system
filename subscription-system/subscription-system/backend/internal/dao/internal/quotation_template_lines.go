// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// QuotationTemplateLinesDao is the data access object for table quotation_template_lines.
type QuotationTemplateLinesDao struct {
	table   string                        // table is the underlying table name of the DAO.
	group   string                        // group is the database configuration group name of current DAO.
	columns QuotationTemplateLinesColumns // columns contains all the column names of Table for convenient usage.
}

// QuotationTemplateLinesColumns defines and stores column names for table quotation_template_lines.
type QuotationTemplateLinesColumns struct {
	Id         string //
	TemplateId string //
	ProductId  string //
	VariantId  string //
	Quantity   string //
}

// quotationTemplateLinesColumns holds the columns for table quotation_template_lines.
var quotationTemplateLinesColumns = QuotationTemplateLinesColumns{
	Id:         "id",
	TemplateId: "template_id",
	ProductId:  "product_id",
	VariantId:  "variant_id",
	Quantity:   "quantity",
}

// NewQuotationTemplateLinesDao creates and returns a new DAO object for table data access.
func NewQuotationTemplateLinesDao() *QuotationTemplateLinesDao {
	return &QuotationTemplateLinesDao{
		group:   "default",
		table:   "quotation_template_lines",
		columns: quotationTemplateLinesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *QuotationTemplateLinesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *QuotationTemplateLinesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *QuotationTemplateLinesDao) Columns() QuotationTemplateLinesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *QuotationTemplateLinesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *QuotationTemplateLinesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *QuotationTemplateLinesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
