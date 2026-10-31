package audit

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"brightloop/api/admin"
	"brightloop/api/common"
)

// List returns audit entries newest first, filtered by table, record, user and an inclusive date range.
func List(ctx context.Context, req *admin.ListAuditReq) (*admin.ListAuditRes, error) {
	m := g.DB().Model("audit_log a").Ctx(ctx).LeftJoin("users u", "u.id = a.changed_by").
		Fields("a.id, a.table_name, a.record_id, a.action, a.changed_by, u.email AS changed_by_email, a.reason, a.changes, a.created_at")
	if req.Table != "" {
		m = m.Where("a.table_name", req.Table)
	}
	if req.RecordID != 0 {
		m = m.Where("a.record_id", req.RecordID)
	}
	if req.UserID != 0 {
		m = m.Where("a.changed_by", req.UserID)
	}
	if req.From != "" {
		m = m.Where("a.created_at >= ?", req.From)
	}
	if req.To != "" {
		m = m.Where("a.created_at < ? + INTERVAL 1 DAY", req.To)
	}
	var items []admin.AuditItem
	var total int
	if err := m.OrderDesc("a.id").Page(req.Page, req.PageSize).ScanAndCount(&items, &total, false); err != nil {
		return nil, err
	}
	return &admin.ListAuditRes{List: common.NewList(items, req.PageReq, total)}, nil
}
