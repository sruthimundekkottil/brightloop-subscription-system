package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"brightloop/internal/logic/auth"
	"brightloop/internal/middleware"
	"brightloop/internal/seed"
	"brightloop/internal/testdb"
)

var base string

func TestMain(m *testing.M) {
	os.Setenv("APP_ENV", "development") // non-Secure cookie over plain http, test clock enabled
	auth.Secret = []byte("test-secret-test-secret-test-secret!")
	if err := testdb.Setup("bl_test_api"); err != nil {
		fmt.Println("testdb:", err)
		os.Exit(1)
	}
	if err := seed.Run(context.Background()); err != nil {
		fmt.Println("seed:", err)
		os.Exit(1)
	}
	s := g.Server("api-test")
	s.SetAddr("127.0.0.1:0")
	s.SetDumpRouterMap(false)
	Routes(s)
	// test-only routes for the role middleware
	s.Group("/api/test", func(grp *ghttp.RouterGroup) {
		grp.Middleware(middleware.Response, middleware.Auth)
		grp.Group("/admin", func(a *ghttp.RouterGroup) {
			a.Middleware(middleware.RequireRole("admin"))
			a.GET("/", func(r *ghttp.Request) { r.Response.WriteJson(g.Map{"data": "ok"}) })
		})
		grp.Group("/back-office", func(a *ghttp.RouterGroup) {
			a.Middleware(middleware.RequireRole("admin", "staff"))
			a.GET("/", func(r *ghttp.Request) { r.Response.WriteJson(g.Map{"data": "ok"}) })
		})
	})
	if err := s.Start(); err != nil {
		fmt.Println("server:", err)
		os.Exit(1)
	}
	base = fmt.Sprintf("http://127.0.0.1:%d/api", s.GetListenedPort())
	code := m.Run()
	s.Shutdown()
	os.Exit(code)
}

// client is one browser: its own cookie jar.
type client struct{ http.Client }

func newClient() *client {
	jar, _ := cookiejar.New(nil)
	return &client{http.Client{Jar: jar}}
}

type reply struct {
	Status int
	Body   string
	Data   map[string]any
	Error  struct {
		Code    string            `json:"code"`
		Message string            `json:"message"`
		Fields  map[string]string `json:"fields"`
	}
}

func (c *client) do(t *testing.T, method, path string, body any) reply {
	t.Helper()
	var rd io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rd = strings.NewReader(string(b))
	}
	req, _ := http.NewRequest(method, base+path, rd)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	r := reply{Status: resp.StatusCode, Body: string(b)}
	var env struct {
		Data  json.RawMessage `json:"data"`
		Error json.RawMessage `json:"error"`
	}
	json.Unmarshal(b, &env)
	json.Unmarshal(env.Data, &r.Data)
	json.Unmarshal(env.Error, &r.Error)
	return r
}

func (c *client) login(t *testing.T, email, password string) reply {
	t.Helper()
	return c.do(t, "POST", "/auth/login", g.Map{"email": email, "password": password})
}

func want(t *testing.T, r reply, status int) {
	t.Helper()
	if r.Status != status {
		t.Fatalf("status %d, want %d: %s", r.Status, status, r.Body)
	}
}

func TestRegister(t *testing.T) {
	c := newClient()
	r := c.do(t, "POST", "/auth/register", g.Map{"email": "weak@shop.test", "password": "short", "name": "Weak"})
	want(t, r, 422)
	if !strings.Contains(r.Error.Fields["password"], "more than 8") || !strings.Contains(r.Error.Fields["password"], "uppercase") {
		t.Errorf("password field message: %q", r.Error.Fields["password"])
	}
	r = c.do(t, "POST", "/auth/register", g.Map{"email": "not-an-email", "password": "Good#Pass1", "name": "X"})
	want(t, r, 422)
	if r.Error.Fields["email"] == "" {
		t.Errorf("no email field error: %s", r.Body)
	}

	r = c.do(t, "POST", "/auth/register", g.Map{"email": "New.Buyer@Shop.test", "password": "Good#Pass1", "name": "New Buyer"})
	want(t, r, 201)
	if r.Data["email"] != "new.buyer@shop.test" || r.Data["role"] != "customer" || r.Data["customerId"] == nil {
		t.Fatalf("register reply: %s", r.Body)
	}
	me := c.do(t, "GET", "/auth/me", nil) // registering logs you in
	want(t, me, 200)
	if me.Data["id"] != r.Data["id"] {
		t.Errorf("me %v after register %v", me.Data, r.Data)
	}
	uid := int(r.Data["id"].(float64))
	if n, _ := g.DB().GetCount(context.Background(), "SELECT COUNT(*) FROM audit_log WHERE changed_by = ? AND action IN ('register', 'create')", uid); n != 2 {
		t.Errorf("want 2 audit rows by the new user, got %d", n)
	}

	r = newClient().do(t, "POST", "/auth/register", g.Map{"email": "NEW.BUYER@shop.TEST", "password": "Good#Pass1", "name": "Dup"})
	want(t, r, 409)
	if r.Error.Fields["email"] != "This email is already registered." {
		t.Errorf("duplicate reply: %s", r.Body)
	}
	if n, _ := g.DB().GetCount(context.Background(), "SELECT COUNT(*) FROM customers WHERE name = 'Dup'"); n != 0 {
		t.Error("failed registration left a customer row")
	}
}

func TestLogin(t *testing.T) {
	c := newClient()
	wrongPw := c.login(t, seed.AdminEmail, "Wrong#Pass1")
	unknown := c.login(t, "nobody@brightloop.test", "Wrong#Pass1")
	want(t, wrongPw, 401)
	if wrongPw.Body != unknown.Body || wrongPw.Status != unknown.Status {
		t.Errorf("wrong password and unknown email differ:\n%s\n%s", wrongPw.Body, unknown.Body)
	}

	r := c.login(t, " ADMIN@brightloop.test", seed.DevPassword)
	want(t, r, 200)
	if r.Data["role"] != "admin" || r.Data["customerId"] != nil {
		t.Errorf("admin login: %s", r.Body)
	}
	want(t, c.do(t, "GET", "/auth/me", nil), 200)
	want(t, c.do(t, "POST", "/auth/logout", nil), 200)
	if r := c.do(t, "GET", "/auth/me", nil); r.Status != 401 || r.Error.Code != "UNAUTHENTICATED" {
		t.Errorf("me after logout: %s", r.Body)
	}
}

func TestDeactivatedUser(t *testing.T) {
	ctx := context.Background()
	c := newClient()
	want(t, c.do(t, "POST", "/auth/register", g.Map{"email": "leaving@shop.test", "password": "Good#Pass1", "name": "Leaving"}), 201)
	g.DB().Exec(ctx, "UPDATE users SET is_active = FALSE WHERE email = 'leaving@shop.test'")
	want(t, c.do(t, "GET", "/auth/me", nil), 401) // the open session dies at once
	want(t, newClient().login(t, "leaving@shop.test", "Good#Pass1"), 403)
	want(t, newClient().login(t, "leaving@shop.test", "Wrong#Pass1"), 401) // wrong password still says nothing more
}

func TestSessionCookieIsRequiredAndSigned(t *testing.T) {
	c := newClient()
	want(t, c.do(t, "GET", "/auth/me", nil), 401)

	want(t, c.login(t, seed.StaffEmail, seed.DevPassword), 200)
	u, _ := http.NewRequest("GET", base, nil)
	var cookie *http.Cookie
	for _, ck := range c.Jar.Cookies(u.URL) {
		if ck.Name == auth.CookieName {
			cookie = ck
		}
	}
	if cookie == nil {
		t.Fatal("no session cookie")
	}
	forged := newClient()
	forged.Jar.SetCookies(u.URL, []*http.Cookie{{Name: auth.CookieName, Value: cookie.Value + "x", Path: "/"}})
	want(t, forged.do(t, "GET", "/auth/me", nil), 401)
	// alg "none" with the admin's id
	forged.Jar.SetCookies(u.URL, []*http.Cookie{{Name: auth.CookieName, Path: "/",
		Value: "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiIxIiwiZXhwIjo5OTk5OTk5OTk5fQ."}})
	want(t, forged.do(t, "GET", "/auth/me", nil), 401)
}

func TestSetCookieFlags(t *testing.T) {
	resp, err := http.Post(base+"/auth/login", "application/json",
		strings.NewReader(fmt.Sprintf(`{"email":%q,"password":%q}`, seed.CustomerEmail, seed.DevPassword)))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	sc := resp.Header.Get("Set-Cookie")
	for _, f := range []string{auth.CookieName + "=", "HttpOnly", "SameSite=Lax", "Path=/"} {
		if !strings.Contains(sc, f) {
			t.Errorf("Set-Cookie %q lacks %s", sc, f)
		}
	}
}

func TestRequireRole(t *testing.T) {
	for _, c := range []struct {
		email             string
		admin, backOffice int
	}{
		{seed.AdminEmail, 200, 200},
		{seed.StaffEmail, 403, 200},
		{seed.CustomerEmail, 403, 403},
	} {
		cl := newClient()
		want(t, cl.login(t, c.email, seed.DevPassword), 200)
		if r := cl.do(t, "GET", "/test/admin/", nil); r.Status != c.admin {
			t.Errorf("%s on admin route: %d, want %d", c.email, r.Status, c.admin)
		}
		if r := cl.do(t, "GET", "/test/back-office/", nil); r.Status != c.backOffice {
			t.Errorf("%s on back-office route: %d, want %d", c.email, r.Status, c.backOffice)
		}
	}
	want(t, newClient().do(t, "GET", "/test/admin/", nil), 401)
}
