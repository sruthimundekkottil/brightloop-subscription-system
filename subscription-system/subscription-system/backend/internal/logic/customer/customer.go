// Package customer manages customer profiles for the back office and the portal.
package customer

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"brightloop/api/admin"
	"brightloop/api/common"
	"brightloop/internal/apperr"
	"brightloop/internal/logic/platform"
	"brightloop/internal/logic/platform/audit"
)

func model(ctx context.Context) *gdb.Model {
	return g.DB().Model("customers c").Ctx(ctx).LeftJoin("users u", "u.id = c.user_id").
		Fields("c.id, c.user_id, u.email, c.name, c.company_name, c.phone, c.billing_address, c.created_at")
}

func List(ctx context.Context, req *admin.ListCustomersReq) (*admin.ListCustomersRes, error) {
	m := model(ctx)
	if req.Q != "" {
		l := platform.Like(req.Q)
		m = m.Where("(c.name LIKE ? OR c.company_name LIKE ? OR u.email LIKE ?)", l, l, l)
	}
	var items []admin.CustomerItem
	var total int
	if err := m.OrderAsc("c.name").OrderAsc("c.id").Page(req.Page, req.PageSize).ScanAndCount(&items, &total, false); err != nil {
		return nil, err
	}
	return &admin.ListCustomersRes{List: common.NewList(items, req.PageReq, total)}, nil
}

// Get returns the customer or NOT_FOUND.
func Get(ctx context.Context, id uint64) (*admin.CustomerItem, error) {
	var c *admin.CustomerItem
	if err := model(ctx).Where("c.id", id).Scan(&c); err != nil {
		return nil, err
	}
	if c == nil {
		return nil, apperr.NotFound()
	}
	return c, nil
}

// columns maps the sent fields to columns; "" clears an optional field, a blank name is rejected.
func columns(f admin.CustomerFields) (map[string]any, error) {
	set := map[string]any{}
	if f.Name != nil {
		if strings.TrimSpace(*f.Name) == "" {
			return nil, apperr.Validation(map[string]string{"name": "Enter a name."})
		}
		set["name"] = strings.TrimSpace(*f.Name)
	}
	for col, v := range map[string]*string{"company_name": f.CompanyName, "phone": f.Phone, "billing_address": f.BillingAddress} {
		if v != nil {
			if s := strings.TrimSpace(*v); s != "" {
				set[col] = s
			} else {
				set[col] = nil
			}
		}
	}
	return set, nil
}

// Create adds a customer without a login (user_id NULL); they can be linked when they register later.
func Create(ctx context.Context, f admin.CustomerFields) (*admin.CustomerItem, error) {
	if f.Name == nil {
		return nil, apperr.Validation(map[string]string{"name": "Enter a name."})
	}
	set, err := columns(f)
	if err != nil {
		return nil, err
	}
	var id int64
	err = platform.Tx(ctx, func(ctx context.Context, tx gdb.TX) error {
		if id, err = tx.Ctx(ctx).Model("customers").Data(set).InsertAndGetId(); err != nil {
			return err
		}
		return audit.Log(ctx, tx, "customers", uint64(id), "create", "", set)
	})
	if err != nil {
		return nil, err
	}
	return Get(ctx, uint64(id))
}

// Update changes the sent fields of customer id. Portal callers must pass auth.CustomerID(ctx) as id.
func Update(ctx context.Context, id uint64, f admin.CustomerFields) (*admin.CustomerItem, error) {
	set, err := columns(f)
	if err != nil {
		return nil, err
	}
	var found bool
	err = platform.Tx(ctx, func(ctx context.Context, tx gdb.TX) error {
		found, err = audit.Update(ctx, tx, "customers", id, set)
		return err
	})
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, apperr.NotFound()
	}
	return Get(ctx, id)
}
