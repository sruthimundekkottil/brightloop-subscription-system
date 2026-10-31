// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CompanySettings is the golang structure for table company_settings.
type CompanySettings struct {
	Id                           uint        `json:"id"                           ` //
	CompanyName                  string      `json:"companyName"                  ` //
	Address                      string      `json:"address"                      ` //
	CurrencyCode                 string      `json:"currencyCode"                 ` //
	InvoicePrefix                string      `json:"invoicePrefix"                ` //
	DefaultQuotationValidityDays uint        `json:"defaultQuotationValidityDays" ` //
	AutoConfirmRecurringInvoices int         `json:"autoConfirmRecurringInvoices" ` //
	TestClockDate                *gtime.Time `json:"testClockDate"                ` //
	UpdatedAt                    *gtime.Time `json:"updatedAt"                    ` //
}
