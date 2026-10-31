package middleware

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"

	"brightloop/api/common"
	"brightloop/internal/apperr"
)

type errReq struct {
	g.Meta `path:"/err" method:"get"`
	Kind   string `json:"kind"`
}
type okRes struct {
	Name string `json:"name"`
}
type createReq struct {
	g.Meta `path:"/things" method:"post"`
	Email  string `json:"email" v:"required|email"`
}
type listReq struct {
	g.Meta `path:"/things" method:"get"`
	common.PageReq
}

type listRes struct{ common.List[okRes] }

type ctl struct{}

func (ctl) Err(ctx context.Context, req *errReq) (*okRes, error) {
	switch req.Kind {
	case "dup":
		return nil, fmt.Errorf("insert: %w", &mysql.MySQLError{Number: 1062, Message: "Duplicate entry 'secret' for key"})
	case "check":
		return nil, &mysql.MySQLError{Number: 3819}
	case "deadlock":
		return nil, &mysql.MySQLError{Number: 1213}
	case "notfound":
		return nil, apperr.NotFound()
	case "raw":
		return nil, errors.New("db password is hunter2")
	case "panic":
		panic("boom")
	}
	return &okRes{Name: "ok"}, nil
}
func (ctl) Create(ctx context.Context, req *createReq) (*okRes, error) { return &okRes{req.Email}, nil }
func (ctl) List(ctx context.Context, req *listReq) (*listRes, error) {
	return &listRes{common.NewList([]okRes{{"a"}}, req.PageReq, 41)}, nil
}

func TestResponseShapes(t *testing.T) {
	s := g.Server(guid.S())
	s.SetAddr("127.0.0.1:0")
	s.SetDumpRouterMap(false)
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(Response)
		group.Bind(ctl{})
	})
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	defer s.Shutdown()
	base := fmt.Sprintf("http://127.0.0.1:%d/api", s.GetListenedPort())

	cases := []struct {
		method, path, body string
		status             int
		want               []string
		notWant            string
	}{
		{"GET", "/err", "", 200, []string{`"data":{"name":"ok"}`}, ""},
		{"GET", "/err?kind=dup", "", 409, []string{`"code":"CONFLICT"`}, "secret"},
		{"GET", "/err?kind=check", "", 422, []string{`"code":"VALIDATION_ERROR"`}, ""},
		{"GET", "/err?kind=deadlock", "", 409, []string{`"code":"CONFLICT"`, "try again"}, ""},
		{"GET", "/err?kind=notfound", "", 404, []string{`"code":"NOT_FOUND"`}, ""},
		{"GET", "/err?kind=raw", "", 500, []string{`"code":"INTERNAL_ERROR"`}, "hunter2"},
		{"GET", "/err?kind=panic", "", 500, []string{`"code":"INTERNAL_ERROR"`}, "boom"},
		{"POST", "/things", `{"email":"nope"}`, 422, []string{`"code":"VALIDATION_ERROR"`, `"fields":{"email":`}, ""},
		{"POST", "/things", `{"email":"a@b.test"}`, 200, []string{`"data":{"name":"a@b.test"}`}, ""},
		{"GET", "/things?page=2&pageSize=20", "", 200, []string{`"data":[{"name":"a"}]`, `"meta":{"page":2,"pageSize":20,"total":41}`}, ""},
		{"GET", "/things", "", 200, []string{`"meta":{"page":1,"pageSize":25,"total":41}`}, ""},
		{"GET", "/things?pageSize=101", "", 422, []string{`"fields":{"pageSize":`}, ""},
	}
	for _, c := range cases {
		req, _ := http.NewRequest(c.method, base+c.path, strings.NewReader(c.body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		body := string(b)
		if resp.StatusCode != c.status {
			t.Errorf("%s %s: status %d, want %d (%s)", c.method, c.path, resp.StatusCode, c.status, body)
		}
		for _, w := range c.want {
			if !strings.Contains(body, w) {
				t.Errorf("%s %s: body %s lacks %s", c.method, c.path, body, w)
			}
		}
		if c.notWant != "" && strings.Contains(body, c.notWant) {
			t.Errorf("%s %s: body leaks %q: %s", c.method, c.path, c.notWant, body)
		}
	}
}
