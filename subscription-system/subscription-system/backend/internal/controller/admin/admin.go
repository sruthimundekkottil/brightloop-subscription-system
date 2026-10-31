// Package admin is the thin HTTP layer for staff, users, customers, settings, the audit log and the test clock.
// Which roles reach each type is decided in cmd.Routes.
package admin

import (
	"context"
	"errors"
	"os"
	"time"

	api "brightloop/api/admin"
	"brightloop/api/portal"
	"brightloop/internal/apperr"
	"brightloop/internal/logic/auth"
	"brightloop/internal/logic/customer"
	"brightloop/internal/logic/platform"
	"brightloop/internal/logic/platform/audit"
	"brightloop/internal/logic/platform/clock"
	"brightloop/internal/logic/platform/settings"
)

// AdminOnly: staff and users, writing settings, the audit log.
type AdminOnly struct{}

// BackOffice: admin and staff.
type BackOffice struct{}

// DevClock is registered only when APP_ENV=development.
type DevClock struct{}

// Portal: customers, always scoped to auth.CustomerID.
type Portal struct{}

func (AdminOnly) CreateStaff(ctx context.Context, req *api.CreateStaffReq) (*api.UserRes, error) {
	u, err := auth.CreateStaff(ctx, req)
	if err != nil {
		return nil, err
	}
	return &api.UserRes{UserItem: *u}, nil
}

func (AdminOnly) ListUsers(ctx context.Context, req *api.ListUsersReq) (*api.ListUsersRes, error) {
	return auth.ListUsers(ctx, req)
}

func (AdminOnly) SetUserStatus(ctx context.Context, req *api.SetUserStatusReq) (*api.UserRes, error) {
	u, err := auth.SetUserStatus(ctx, req.ID, *req.IsActive)
	if err != nil {
		return nil, err
	}
	return &api.UserRes{UserItem: *u}, nil
}

func (AdminOnly) UpdateSettings(ctx context.Context, req *api.UpdateSettingsReq) (*api.SettingsRes, error) {
	set := map[string]any{}
	if req.CompanyName != nil {
		set["company_name"] = *req.CompanyName
	}
	if req.Address != nil {
		set["address"] = nilIfEmpty(*req.Address)
	}
	if req.CurrencyCode != nil {
		set["currency_code"] = *req.CurrencyCode
	}
	if req.InvoicePrefix != nil {
		set["invoice_prefix"] = *req.InvoicePrefix
	}
	if req.DefaultQuotationValidityDays != nil {
		set["default_quotation_validity_days"] = *req.DefaultQuotationValidityDays
	}
	if req.AutoConfirmRecurringInvoices != nil {
		set["auto_confirm_recurring_invoices"] = *req.AutoConfirmRecurringInvoices
	}
	if err := settings.Update(ctx, set); err != nil {
		return nil, err
	}
	return settingsRes(ctx)
}

func (AdminOnly) ListAudit(ctx context.Context, req *api.ListAuditReq) (*api.ListAuditRes, error) {
	return audit.List(ctx, req)
}

func (BackOffice) GetSettings(ctx context.Context, _ *api.GetSettingsReq) (*api.SettingsRes, error) {
	return settingsRes(ctx)
}

func (BackOffice) ListCustomers(ctx context.Context, req *api.ListCustomersReq) (*api.ListCustomersRes, error) {
	return customer.List(ctx, req)
}

func (BackOffice) CreateCustomer(ctx context.Context, req *api.CreateCustomerReq) (*api.CustomerRes, error) {
	return customerRes(customer.Create(ctx, req.CustomerFields))
}

func (BackOffice) GetCustomer(ctx context.Context, req *api.GetCustomerReq) (*api.CustomerRes, error) {
	return customerRes(customer.Get(ctx, req.ID))
}

func (BackOffice) UpdateCustomer(ctx context.Context, req *api.UpdateCustomerReq) (*api.CustomerRes, error) {
	return customerRes(customer.Update(ctx, req.ID, req.CustomerFields))
}

func (Portal) GetProfile(ctx context.Context, _ *portal.GetProfileReq) (*api.CustomerRes, error) {
	return customerRes(customer.Get(ctx, auth.CustomerID(ctx)))
}

func (Portal) UpdateProfile(ctx context.Context, req *portal.UpdateProfileReq) (*api.CustomerRes, error) {
	return customerRes(customer.Update(ctx, auth.CustomerID(ctx), req.CustomerFields))
}

func (DevClock) SetClock(ctx context.Context, req *api.SetClockReq) (*api.SettingsRes, error) {
	d, err := time.Parse(platform.DateLayout, req.Date)
	if err != nil {
		return nil, apperr.Validation(map[string]string{"date": "Use YYYY-MM-DD."})
	}
	if err = settings.SetTestClock(ctx, d); errors.Is(err, platform.ErrClockBackwards) {
		return nil, apperr.Conflict("The test clock can only move forward.")
	} else if err != nil {
		return nil, err
	}
	return settingsRes(ctx)
}

func customerRes(c *api.CustomerItem, err error) (*api.CustomerRes, error) {
	if err != nil {
		return nil, err
	}
	return &api.CustomerRes{CustomerItem: *c}, nil
}

func settingsRes(ctx context.Context) (*api.SettingsRes, error) {
	s, err := settings.Get(ctx)
	if err != nil {
		return nil, err
	}
	today, err := clock.Today(ctx)
	if err != nil {
		return nil, err
	}
	res := &api.SettingsRes{
		CompanyName: s.CompanyName, CurrencyCode: s.CurrencyCode, InvoicePrefix: s.InvoicePrefix,
		DefaultQuotationValidityDays: int(s.DefaultQuotationValidityDays), AutoConfirmRecurringInvoices: s.AutoConfirmRecurringInvoices == 1,
		Today: today.Format(platform.DateLayout), DevMode: os.Getenv("APP_ENV") == "development",
	}
	if s.Address != "" {
		res.Address = &s.Address
	}
	if s.TestClockDate != nil {
		d := s.TestClockDate.Format("Y-m-d")
		res.TestClockDate = &d
	}
	return res, nil
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
