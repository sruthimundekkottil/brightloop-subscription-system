// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CompanySettingsDao is the data access object for table company_settings.
type CompanySettingsDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of current DAO.
	columns CompanySettingsColumns // columns contains all the column names of Table for convenient usage.
}

// CompanySettingsColumns defines and stores column names for table company_settings.
type CompanySettingsColumns struct {
	Id                           string //
	CompanyName                  string //
	Address                      string //
	CurrencyCode                 string //
	InvoicePrefix                string //
	DefaultQuotationValidityDays string //
	AutoConfirmRecurringInvoices string //
	TestClockDate                string //
	UpdatedAt                    string //
}

// companySettingsColumns holds the columns for table company_settings.
var companySettingsColumns = CompanySettingsColumns{
	Id:                           "id",
	CompanyName:                  "company_name",
	Address:                      "address",
	CurrencyCode:                 "currency_code",
	InvoicePrefix:                "invoice_prefix",
	DefaultQuotationValidityDays: "default_quotation_validity_days",
	AutoConfirmRecurringInvoices: "auto_confirm_recurring_invoices",
	TestClockDate:                "test_clock_date",
	UpdatedAt:                    "updated_at",
}

// NewCompanySettingsDao creates and returns a new DAO object for table data access.
func NewCompanySettingsDao() *CompanySettingsDao {
	return &CompanySettingsDao{
		group:   "default",
		table:   "company_settings",
		columns: companySettingsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CompanySettingsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CompanySettingsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CompanySettingsDao) Columns() CompanySettingsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CompanySettingsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CompanySettingsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CompanySettingsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
