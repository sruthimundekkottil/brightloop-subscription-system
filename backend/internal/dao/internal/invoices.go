// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InvoicesDao is the data access object for table invoices.
type InvoicesDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns InvoicesColumns // columns contains all the column names of Table for convenient usage.
}

// InvoicesColumns defines and stores column names for table invoices.
type InvoicesColumns struct {
	Id             string //
	Number         string //
	CustomerId     string //
	SubscriptionId string //
	Status         string //
	InvoiceDate    string //
	DueDate        string //
	PeriodStart    string //
	PeriodEnd      string //
	BillingName    string //
	BillingAddress string //
	Subtotal       string //
	DiscountTotal  string //
	TaxTotal       string //
	Total          string //
	AmountPaid     string //
	PaymentStatus  string //
	ConfirmedAt    string //
	SentAt         string //
	CancelledAt    string //
	CancelReason   string //
	CreatedAt      string //
	UpdatedAt      string //
}

// invoicesColumns holds the columns for table invoices.
var invoicesColumns = InvoicesColumns{
	Id:             "id",
	Number:         "number",
	CustomerId:     "customer_id",
	SubscriptionId: "subscription_id",
	Status:         "status",
	InvoiceDate:    "invoice_date",
	DueDate:        "due_date",
	PeriodStart:    "period_start",
	PeriodEnd:      "period_end",
	BillingName:    "billing_name",
	BillingAddress: "billing_address",
	Subtotal:       "subtotal",
	DiscountTotal:  "discount_total",
	TaxTotal:       "tax_total",
	Total:          "total",
	AmountPaid:     "amount_paid",
	PaymentStatus:  "payment_status",
	ConfirmedAt:    "confirmed_at",
	SentAt:         "sent_at",
	CancelledAt:    "cancelled_at",
	CancelReason:   "cancel_reason",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewInvoicesDao creates and returns a new DAO object for table data access.
func NewInvoicesDao() *InvoicesDao {
	return &InvoicesDao{
		group:   "default",
		table:   "invoices",
		columns: invoicesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *InvoicesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *InvoicesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *InvoicesDao) Columns() InvoicesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *InvoicesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *InvoicesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *InvoicesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
