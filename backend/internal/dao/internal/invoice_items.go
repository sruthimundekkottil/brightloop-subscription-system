// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InvoiceItemsDao is the data access object for table invoice_items.
type InvoiceItemsDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns InvoiceItemsColumns // columns contains all the column names of Table for convenient usage.
}

// InvoiceItemsColumns defines and stores column names for table invoice_items.
type InvoiceItemsColumns struct {
	Id             string //
	InvoiceId      string //
	ItemType       string //
	LinkedItemId   string //
	ProductId      string //
	PlanId         string //
	DiscountRuleId string //
	TaxRuleId      string //
	Description    string //
	ServiceStart   string //
	ServiceEnd     string //
	Quantity       string //
	UnitPrice      string //
	RatePercent    string //
	Amount         string //
	CreatedAt      string //
}

// invoiceItemsColumns holds the columns for table invoice_items.
var invoiceItemsColumns = InvoiceItemsColumns{
	Id:             "id",
	InvoiceId:      "invoice_id",
	ItemType:       "item_type",
	LinkedItemId:   "linked_item_id",
	ProductId:      "product_id",
	PlanId:         "plan_id",
	DiscountRuleId: "discount_rule_id",
	TaxRuleId:      "tax_rule_id",
	Description:    "description",
	ServiceStart:   "service_start",
	ServiceEnd:     "service_end",
	Quantity:       "quantity",
	UnitPrice:      "unit_price",
	RatePercent:    "rate_percent",
	Amount:         "amount",
	CreatedAt:      "created_at",
}

// NewInvoiceItemsDao creates and returns a new DAO object for table data access.
func NewInvoiceItemsDao() *InvoiceItemsDao {
	return &InvoiceItemsDao{
		group:   "default",
		table:   "invoice_items",
		columns: invoiceItemsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *InvoiceItemsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *InvoiceItemsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *InvoiceItemsDao) Columns() InvoiceItemsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *InvoiceItemsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *InvoiceItemsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *InvoiceItemsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
