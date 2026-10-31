// Package admin is the request and response shapes of the back-office API (admin and staff).
package admin

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"brightloop/api/common"
)

// ---- users and staff (admin) ----

type CreateStaffReq struct {
	g.Meta   `path:"/staff" method:"post" tags:"Staff" summary:"Create a staff login and profile"`
	Email    string `json:"email"    v:"required|email|max-length:254"`
	Password string `json:"password" v:"required"`
	FullName string `json:"fullName" v:"required|length:1,150"`
	JobTitle string `json:"jobTitle" v:"max-length:100"`
}

type UserItem struct {
	ID         uint64      `json:"id"`
	Email      string      `json:"email"`
	Role       string      `json:"role"`
	IsActive   bool        `json:"isActive"`
	FullName   *string     `json:"fullName"` // staff profile or customer name
	JobTitle   *string     `json:"jobTitle"`
	CustomerID *uint64     `json:"customerId"`
	CreatedAt  *gtime.Time `json:"createdAt"`
}

type UserRes struct{ UserItem }

type ListUsersReq struct {
	g.Meta `path:"/users" method:"get" tags:"Users" summary:"List users"`
	common.PageReq
	Role   string `json:"role"   v:"in:admin,staff,customer"`
	Active string `json:"active" v:"in:true,false"`
	Q      string `json:"q"      v:"max-length:100"` // email or name contains
}
type ListUsersRes struct{ common.List[UserItem] }

type SetUserStatusReq struct {
	g.Meta   `path:"/users/{id}/status" method:"patch" tags:"Users" summary:"Activate or deactivate a user"`
	ID       uint64 `json:"id"       in:"path" v:"required"`
	IsActive *bool  `json:"isActive" v:"required"`
}

// ---- customers (admin, staff) ----

type CustomerItem struct {
	ID             uint64      `json:"id"`
	UserID         *uint64     `json:"userId"` // null: no login yet
	Email          *string     `json:"email"`
	Name           string      `json:"name"`
	CompanyName    *string     `json:"companyName"`
	Phone          *string     `json:"phone"`
	BillingAddress *string     `json:"billingAddress"`
	CreatedAt      *gtime.Time `json:"createdAt"`
}

type CustomerRes struct{ CustomerItem }

type ListCustomersReq struct {
	g.Meta `path:"/customers" method:"get" tags:"Customers" summary:"List customers"`
	common.PageReq
	Q string `json:"q" v:"max-length:100"` // name, company or email contains
}
type ListCustomersRes struct{ common.List[CustomerItem] }

// CustomerFields are the editable fields; nil means "leave as is", "" clears an optional field.
type CustomerFields struct {
	Name           *string `json:"name"           v:"length:1,150"`
	CompanyName    *string `json:"companyName"    v:"max-length:150"`
	Phone          *string `json:"phone"          v:"max-length:30"`
	BillingAddress *string `json:"billingAddress" v:"max-length:1000"`
}

type CreateCustomerReq struct {
	g.Meta `path:"/customers" method:"post" tags:"Customers" summary:"Create a customer without a login"`
	CustomerFields
}

type GetCustomerReq struct {
	g.Meta `path:"/customers/{id}" method:"get" tags:"Customers" summary:"Customer detail"`
	ID     uint64 `json:"id" in:"path" v:"required"`
}

type UpdateCustomerReq struct {
	g.Meta `path:"/customers/{id}" method:"patch" tags:"Customers" summary:"Update a customer"`
	ID     uint64 `json:"id" in:"path" v:"required"`
	CustomerFields
}

// ---- company settings ----

type SettingsRes struct {
	CompanyName                  string  `json:"companyName"`
	Address                      *string `json:"address"`
	CurrencyCode                 string  `json:"currencyCode"`
	InvoicePrefix                string  `json:"invoicePrefix"`
	DefaultQuotationValidityDays int     `json:"defaultQuotationValidityDays"`
	AutoConfirmRecurringInvoices bool    `json:"autoConfirmRecurringInvoices"`
	TestClockDate                *string `json:"testClockDate"` // YYYY-MM-DD or null
	Today                        string  `json:"today"`         // the business date every module uses
	DevMode                      bool    `json:"devMode"`       // the test clock endpoint exists
}

type GetSettingsReq struct {
	g.Meta `path:"/settings" method:"get" tags:"Settings" summary:"Company settings"`
}

// UpdateSettingsReq changes only the fields sent.
type UpdateSettingsReq struct {
	g.Meta                       `path:"/settings" method:"put" tags:"Settings" summary:"Update company settings"`
	CompanyName                  *string `json:"companyName"                  v:"length:1,150"`
	Address                      *string `json:"address"                      v:"max-length:1000"`
	CurrencyCode                 *string `json:"currencyCode"                 v:"regex:^[A-Z]{3}$"`
	InvoicePrefix                *string `json:"invoicePrefix"                v:"length:1,10|regex:^[A-Z0-9]+$"`
	DefaultQuotationValidityDays *int    `json:"defaultQuotationValidityDays" v:"between:1,365"`
	AutoConfirmRecurringInvoices *bool   `json:"autoConfirmRecurringInvoices"`
}

type SetClockReq struct {
	g.Meta `path:"/dev/clock" method:"put" tags:"Development" summary:"Move the test clock forward (development only)"`
	Date   string `json:"date" v:"required|date-format:Y-m-d"`
}

// ---- audit log (admin) ----

type AuditItem struct {
	ID             uint64      `json:"id"`
	TableName      string      `json:"tableName"`
	RecordID       uint64      `json:"recordId"`
	Action         string      `json:"action"`
	ChangedBy      *uint64     `json:"changedBy"` // null: system job
	ChangedByEmail *string     `json:"changedByEmail"`
	Reason         *string     `json:"reason"`
	Changes        *gjson.Json `json:"changes"`
	CreatedAt      *gtime.Time `json:"createdAt"`
}

type ListAuditReq struct {
	g.Meta `path:"/audit-log" method:"get" tags:"Audit" summary:"Audit log, newest first"`
	common.PageReq
	Table    string `json:"table"    v:"max-length:64"`
	RecordID uint64 `json:"recordId"`
	UserID   uint64 `json:"userId"`
	From     string `json:"from"     v:"date-format:Y-m-d"`
	To       string `json:"to"       v:"date-format:Y-m-d"` // inclusive
}
type ListAuditRes struct{ common.List[AuditItem] }
