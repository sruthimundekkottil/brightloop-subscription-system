// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SubscriptionsDao is the data access object for table subscriptions.
type SubscriptionsDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns SubscriptionsColumns // columns contains all the column names of Table for convenient usage.
}

// SubscriptionsColumns defines and stores column names for table subscriptions.
type SubscriptionsColumns struct {
	Id                  string //
	Number              string //
	CustomerId          string //
	PlanId              string //
	TemplateId          string //
	PaymentTermsId      string //
	Status              string //
	StartDate           string //
	ExpiryDate          string //
	BillingAnchorDate   string //
	NextBillingDate     string //
	QuotationValidUntil string //
	ScheduledCloseDate  string //
	RecurringTotal      string //
	RecurringMonthly    string //
	RenewedFromId       string //
	Terms               string //
	SalespersonId       string //
	PausedAt            string //
	ClosedAt            string //
	CloseReasonId       string //
	CloseNote           string //
	CreatedAt           string //
	UpdatedAt           string //
	OpenRenewalMarker   string //
}

// subscriptionsColumns holds the columns for table subscriptions.
var subscriptionsColumns = SubscriptionsColumns{
	Id:                  "id",
	Number:              "number",
	CustomerId:          "customer_id",
	PlanId:              "plan_id",
	TemplateId:          "template_id",
	PaymentTermsId:      "payment_terms_id",
	Status:              "status",
	StartDate:           "start_date",
	ExpiryDate:          "expiry_date",
	BillingAnchorDate:   "billing_anchor_date",
	NextBillingDate:     "next_billing_date",
	QuotationValidUntil: "quotation_valid_until",
	ScheduledCloseDate:  "scheduled_close_date",
	RecurringTotal:      "recurring_total",
	RecurringMonthly:    "recurring_monthly",
	RenewedFromId:       "renewed_from_id",
	Terms:               "terms",
	SalespersonId:       "salesperson_id",
	PausedAt:            "paused_at",
	ClosedAt:            "closed_at",
	CloseReasonId:       "close_reason_id",
	CloseNote:           "close_note",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
	OpenRenewalMarker:   "open_renewal_marker",
}

// NewSubscriptionsDao creates and returns a new DAO object for table data access.
func NewSubscriptionsDao() *SubscriptionsDao {
	return &SubscriptionsDao{
		group:   "default",
		table:   "subscriptions",
		columns: subscriptionsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SubscriptionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SubscriptionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SubscriptionsDao) Columns() SubscriptionsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SubscriptionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SubscriptionsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SubscriptionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
