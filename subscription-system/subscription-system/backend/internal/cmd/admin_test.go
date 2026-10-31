package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"brightloop/internal/seed"
)

func as(t *testing.T, email string) *client {
	t.Helper()
	c := newClient()
	want(t, c.login(t, email, seed.DevPassword), 200)
	return c
}

func lastAudit(t *testing.T, table string, id any) g.Map {
	t.Helper()
	r, _ := g.DB().GetOne(context.Background(), `SELECT action, changed_by, JSON_UNQUOTE(changes) AS changes FROM audit_log
		WHERE table_name = ? AND record_id = ? ORDER BY id DESC LIMIT 1`, table, id)
	return r.Map()
}

func adminID() int {
	v, _ := g.DB().GetValue(context.Background(), "SELECT id FROM users WHERE email = ?", seed.AdminEmail)
	return v.Int()
}

func TestStaffAndUsers(t *testing.T) {
	adm, staff, cust := as(t, seed.AdminEmail), as(t, seed.StaffEmail), as(t, seed.CustomerEmail)
	body := g.Map{"email": "new.staff@brightloop.test", "password": "Staff#Pass1", "fullName": "Nia Staff", "jobTitle": "Support"}

	want(t, staff.do(t, "POST", "/staff", body), 403)
	want(t, cust.do(t, "POST", "/staff", body), 403)
	want(t, newClient().do(t, "POST", "/staff", body), 401)
	weak := g.Map{"email": "weak.staff@brightloop.test", "password": "weak", "fullName": "W"}
	if r := adm.do(t, "POST", "/staff", weak); r.Status != 422 || r.Error.Fields["password"] == "" {
		t.Errorf("weak staff password: %s", r.Body)
	}

	r := adm.do(t, "POST", "/staff", body)
	want(t, r, 200)
	if r.Data["role"] != "staff" || r.Data["fullName"] != "Nia Staff" || r.Data["isActive"] != true {
		t.Fatalf("create staff: %s", r.Body)
	}
	id := int(r.Data["id"].(float64))
	if a := lastAudit(t, "users", id); a["action"] != "create-staff" || g.NewVar(a["changed_by"]).Int() != adminID() {
		t.Errorf("staff audit: %v", a)
	}
	if by, _ := g.DB().GetValue(context.Background(), "SELECT created_by FROM staff_profiles WHERE user_id = ?", id); by.Int() != adminID() {
		t.Error("staff_profiles.created_by is not the admin")
	}
	want(t, adm.do(t, "POST", "/staff", g.Map{"email": "NEW.STAFF@brightloop.test", "password": "Staff#Pass1", "fullName": "Dup"}), 409)
	newStaff := newClient()
	want(t, newStaff.login(t, "new.staff@brightloop.test", "Staff#Pass1"), 200)

	r = adm.do(t, "GET", "/users?role=staff&pageSize=100", nil)
	want(t, r, 200)
	if r.Meta.Total < 2 || len(r.List) != r.Meta.Total || r.Meta.PageSize != 100 {
		t.Errorf("staff list: %s", r.Body)
	}
	for _, u := range r.List {
		if u["role"] != "staff" {
			t.Errorf("role filter leaked %v", u)
		}
	}
	if r = adm.do(t, "GET", "/users?q=nia", nil); len(r.List) != 1 {
		t.Errorf("search by name: %s", r.Body)
	}
	want(t, adm.do(t, "GET", "/users?pageSize=101", nil), 422)
	want(t, staff.do(t, "GET", "/users", nil), 403)

	if r = adm.do(t, "PATCH", fmt.Sprintf("/users/%d/status", adminID()), g.Map{"isActive": false}); r.Status != 409 {
		t.Errorf("self-deactivation: %s", r.Body)
	}
	want(t, adm.do(t, "PATCH", "/users/999999/status", g.Map{"isActive": false}), 404)
	want(t, adm.do(t, "PATCH", fmt.Sprintf("/users/%d/status", id), g.Map{}), 422)
	r = adm.do(t, "PATCH", fmt.Sprintf("/users/%d/status", id), g.Map{"isActive": false})
	want(t, r, 200)
	if r.Data["isActive"] != false {
		t.Errorf("deactivate: %s", r.Body)
	}
	if a := lastAudit(t, "users", id); a["action"] != "update" || a["changes"] != `{"is_active": [1, 0]}` {
		t.Errorf("status audit: %v", a)
	}
	want(t, newStaff.do(t, "GET", "/auth/me", nil), 401)
	want(t, adm.do(t, "PATCH", fmt.Sprintf("/users/%d/status", id), g.Map{"isActive": true}), 200)
}

func TestCustomers(t *testing.T) {
	staff, cust := as(t, seed.StaffEmail), as(t, seed.CustomerEmail)
	want(t, cust.do(t, "GET", "/customers", nil), 403)

	if r := staff.do(t, "POST", "/customers", g.Map{"companyName": "No Name Ltd"}); r.Status != 422 || r.Error.Fields["name"] == "" {
		t.Errorf("missing name: %s", r.Body)
	}
	r := staff.do(t, "POST", "/customers", g.Map{"name": "Walk-in Ltd", "companyName": "Walk-in", "phone": "+91 98470 00000"})
	want(t, r, 200)
	if r.Data["userId"] != nil || r.Data["email"] != nil {
		t.Errorf("staff-created customer should have no login: %s", r.Body)
	}
	id := int(r.Data["id"].(float64))

	r = staff.do(t, "GET", "/customers?q=walk-in", nil)
	if len(r.List) != 1 || r.Meta.Total != 1 {
		t.Errorf("search: %s", r.Body)
	}
	if r = staff.do(t, "GET", "/customers?q=customer@brightloop", nil); len(r.List) != 1 {
		t.Errorf("search by login email: %s", r.Body)
	}
	if r = staff.do(t, "GET", "/customers?q=%25", nil); r.Meta.Total != 0 {
		t.Errorf("%% must match literally, got %d", r.Meta.Total)
	}

	want(t, staff.do(t, "PATCH", fmt.Sprintf("/customers/%d", id), g.Map{"name": "  "}), 422)
	r = staff.do(t, "PATCH", fmt.Sprintf("/customers/%d", id), g.Map{"name": "Walk-in Pvt Ltd", "companyName": ""})
	want(t, r, 200)
	if r.Data["name"] != "Walk-in Pvt Ltd" || r.Data["companyName"] != nil || r.Data["phone"] != "+91 98470 00000" {
		t.Errorf("patch: %s", r.Body)
	}
	if a := lastAudit(t, "customers", id); a["changes"] != `{"name": ["Walk-in Ltd", "Walk-in Pvt Ltd"], "company_name": ["Walk-in", null]}` {
		t.Errorf("customer audit: %v", a)
	}
	auditRows := func() int {
		n, _ := g.DB().GetCount(context.Background(), "SELECT COUNT(*) FROM audit_log WHERE table_name = 'customers' AND record_id = ?", id)
		return n
	}
	before := auditRows()
	want(t, staff.do(t, "PATCH", fmt.Sprintf("/customers/%d", id), g.Map{"name": "Walk-in Pvt Ltd"}), 200)
	if auditRows() != before {
		t.Error("a no-op update wrote an audit row")
	}
	want(t, staff.do(t, "GET", "/customers/999999", nil), 404)
	want(t, staff.do(t, "PATCH", "/customers/999999", g.Map{"name": "X"}), 404)
}

func TestPortalProfileIsOwnOnly(t *testing.T) {
	a, b := newClient(), newClient()
	want(t, a.do(t, "POST", "/auth/register", g.Map{"email": "a@portal.test", "password": "Good#Pass1", "name": "Alpha"}), 201)
	want(t, b.do(t, "POST", "/auth/register", g.Map{"email": "b@portal.test", "password": "Good#Pass1", "name": "Beta"}), 201)

	r := a.do(t, "GET", "/portal/profile", nil)
	if r.Status != 200 || r.Data["name"] != "Alpha" {
		t.Fatalf("own profile: %s", r.Body)
	}
	bID := int(b.do(t, "GET", "/portal/profile", nil).Data["id"].(float64))
	want(t, a.do(t, "GET", fmt.Sprintf("/customers/%d", bID), nil), 403) // back-office route is closed to customers
	want(t, a.do(t, "PATCH", "/portal/profile", g.Map{"phone": "12345"}), 200)
	if r = b.do(t, "GET", "/portal/profile", nil); r.Data["phone"] != nil || r.Data["name"] != "Beta" {
		t.Errorf("A's update touched B: %s", r.Body)
	}
	want(t, as(t, seed.StaffEmail).do(t, "GET", "/portal/profile", nil), 403)
}

func TestSettings(t *testing.T) {
	adm, staff := as(t, seed.AdminEmail), as(t, seed.StaffEmail)
	r := staff.do(t, "GET", "/settings", nil)
	if r.Status != 200 || r.Data["invoicePrefix"] != "INV" || r.Data["devMode"] != true || r.Data["today"] == nil {
		t.Fatalf("settings: %s", r.Body)
	}
	want(t, staff.do(t, "PUT", "/settings", g.Map{"companyName": "Hacked"}), 403)
	// GoFrame validates requests with Bail(): one field error per reply.
	for field, bad := range map[string]any{"invoicePrefix": "inv-", "currencyCode": "rupees", "defaultQuotationValidityDays": 0} {
		if r = adm.do(t, "PUT", "/settings", g.Map{field: bad}); r.Status != 422 || r.Error.Fields[field] == "" {
			t.Errorf("bad %s: %s", field, r.Body)
		}
	}
	if r = adm.do(t, "PUT", "/settings", g.Map{"invoicePrefix": "TOOLONGPREFIX"}); r.Status != 422 {
		t.Errorf("11-char prefix accepted: %s", r.Body)
	}
	r = adm.do(t, "PUT", "/settings", g.Map{"companyName": "Brightloop Labs", "autoConfirmRecurringInvoices": true})
	if r.Status != 200 || r.Data["companyName"] != "Brightloop Labs" || r.Data["autoConfirmRecurringInvoices"] != true || r.Data["invoicePrefix"] != "INV" {
		t.Errorf("update settings: %s", r.Body)
	}
	if a := lastAudit(t, "company_settings", 1); a["action"] != "update" || g.NewVar(a["changed_by"]).Int() != adminID() {
		t.Errorf("settings audit: %v", a)
	}
	adm.do(t, "PUT", "/settings", g.Map{"companyName": "Brightloop Software", "autoConfirmRecurringInvoices": false})
}

func TestAuditLogAPI(t *testing.T) {
	adm := as(t, seed.AdminEmail)
	want(t, as(t, seed.StaffEmail).do(t, "GET", "/audit-log", nil), 403)
	adm.do(t, "PUT", "/settings", g.Map{"address": "1 Audit Street"})
	adm.do(t, "PUT", "/settings", g.Map{"address": ""})
	r := adm.do(t, "GET", fmt.Sprintf("/audit-log?table=company_settings&recordId=1&userId=%d&pageSize=2", adminID()), nil)
	want(t, r, 200)
	if len(r.List) != 2 || r.List[0]["changedByEmail"] != seed.AdminEmail {
		t.Fatalf("audit list: %s", r.Body)
	}
	if ch, ok := r.List[0]["changes"].(map[string]any); !ok || ch["address"] == nil { // newest first, changes as JSON
		t.Errorf("changes not an object or wrong order: %v", r.List[0])
	}
	today := time.Now().UTC().Format("2006-01-02")
	if r = adm.do(t, "GET", "/audit-log?from="+today+"&to="+today+"&table=company_settings", nil); r.Meta.Total < 2 {
		t.Errorf("date range: %s", r.Body)
	}
	if r = adm.do(t, "GET", "/audit-log?to=2000-01-01", nil); r.Meta.Total != 0 {
		t.Errorf("to filter: %s", r.Body)
	}
	want(t, adm.do(t, "GET", "/audit-log?from=31-12-2026", nil), 422)
}

func TestDevClock(t *testing.T) {
	adm := as(t, seed.AdminEmail)
	ahead := time.Now().UTC().AddDate(0, 0, 10).Format("2006-01-02")
	behind := time.Now().UTC().AddDate(0, 0, 5).Format("2006-01-02")

	want(t, as(t, seed.StaffEmail).do(t, "PUT", "/dev/clock", g.Map{"date": ahead}), 403)
	want(t, adm.do(t, "PUT", "/dev/clock", g.Map{"date": "tomorrow"}), 422)
	r := adm.do(t, "PUT", "/dev/clock", g.Map{"date": ahead})
	if r.Status != 200 || r.Data["today"] != ahead || r.Data["testClockDate"] != ahead {
		t.Fatalf("set clock: %s", r.Body)
	}
	if a := lastAudit(t, "company_settings", 1); a["action"] != "test-clock" {
		t.Errorf("clock audit: %v", a)
	}
	if r = adm.do(t, "PUT", "/dev/clock", g.Map{"date": behind}); r.Status != 409 {
		t.Errorf("backwards: %s", r.Body)
	}
	want(t, adm.do(t, "PUT", "/dev/clock", g.Map{"date": ahead}), 409) // same day is not forward

	// Outside development the route does not exist at all.
	os.Setenv("APP_ENV", "production")
	s := g.Server("api-test-prod")
	s.SetAddr("127.0.0.1:0")
	s.SetDumpRouterMap(false)
	Routes(s)
	os.Setenv("APP_ENV", "development")
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	defer s.Shutdown()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("http://127.0.0.1:%d/api/dev/clock", s.GetListenedPort()), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Errorf("production /api/dev/clock: %d, want 404", resp.StatusCode)
	}
	want(t, newClient().do(t, "PUT", "/dev/clock", g.Map{"date": ahead}), 401) // exists in development
}
