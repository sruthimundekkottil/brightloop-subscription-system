package auth

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

func usersModel(ctx context.Context) *gdb.Model {
	return g.DB().Model("users u").Ctx(ctx).
		LeftJoin("staff_profiles s", "s.user_id = u.id").
		LeftJoin("customers c", "c.user_id = u.id").
		Fields("u.id, u.email, u.role, u.is_active, COALESCE(s.full_name, c.name) AS full_name, s.job_title, c.id AS customer_id, u.created_at")
}

// CreateStaff creates a staff login and its profile in one transaction. The caller must be an admin (route guard).
func CreateStaff(ctx context.Context, req *admin.CreateStaffReq) (*admin.UserItem, error) {
	hash, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	var uid int64
	err = platform.Tx(ctx, func(ctx context.Context, tx gdb.TX) error {
		email := NormalizeEmail(req.Email)
		r, err := tx.Exec("INSERT INTO users (email, password_hash, role) VALUES (?, ?, 'staff')", email, hash)
		if err != nil {
			return err
		}
		uid, _ = r.LastInsertId()
		if _, err = tx.Exec("INSERT INTO staff_profiles (user_id, full_name, job_title, created_by) VALUES (?, ?, NULLIF(?, ''), ?)",
			uid, strings.TrimSpace(req.FullName), strings.TrimSpace(req.JobTitle), platform.ActorFrom(ctx).UserID); err != nil {
			return err
		}
		return audit.Log(ctx, tx, "users", uint64(uid), "create-staff", "", g.Map{"email": email, "fullName": req.FullName})
	})
	if IsDuplicate(err) {
		return nil, EmailTaken()
	}
	if err != nil {
		return nil, err
	}
	var item *admin.UserItem
	return item, usersModel(ctx).Where("u.id", uid).Scan(&item)
}

func ListUsers(ctx context.Context, req *admin.ListUsersReq) (*admin.ListUsersRes, error) {
	m := usersModel(ctx)
	if req.Role != "" {
		m = m.Where("u.role", req.Role)
	}
	if req.Active != "" {
		m = m.Where("u.is_active", req.Active == "true")
	}
	if req.Q != "" {
		l := platform.Like(req.Q)
		m = m.Where("(u.email LIKE ? OR s.full_name LIKE ? OR c.name LIKE ?)", l, l, l)
	}
	var items []admin.UserItem
	var total int
	if err := m.OrderDesc("u.id").Page(req.Page, req.PageSize).ScanAndCount(&items, &total, false); err != nil {
		return nil, err
	}
	return &admin.ListUsersRes{List: common.NewList(items, req.PageReq, total)}, nil
}

// SetUserStatus activates or deactivates a user; an admin cannot deactivate themselves.
// A deactivated user's open sessions stop working at once (middleware.Auth re-reads is_active).
func SetUserStatus(ctx context.Context, id uint64, active bool) (*admin.UserItem, error) {
	if id == platform.ActorFrom(ctx).UserID && !active {
		return nil, apperr.Conflict("You cannot deactivate your own account.")
	}
	var found bool
	err := platform.Tx(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		found, err = audit.Update(ctx, tx, "users", id, map[string]any{"is_active": active})
		return err
	})
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, apperr.NotFound()
	}
	var item *admin.UserItem
	return item, usersModel(ctx).Where("u.id", id).Scan(&item)
}
