// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CompanySettings is the golang structure of table company_settings for DAO operations like Where/Data.
type CompanySettings struct {
	g.Meta                       `orm:"table:company_settings, do:true"`
	Id                           interface{} //
	CompanyName                  interface{} //
	Address                      interface{} //
	CurrencyCode                 interface{} //
	InvoicePrefix                interface{} //
	DefaultQuotationValidityDays interface{} //
	AutoConfirmRecurringInvoices interface{} //
	TestClockDate                *gtime.Time //
	UpdatedAt                    *gtime.Time //
}
