// Package cmd wires the HTTP routes. main and the API tests both use Routes, so tests exercise the real wiring.
package cmd

import (
	"os"

	"github.com/gogf/gf/v2/net/ghttp"

	adminctl "brightloop/internal/controller/admin"
	authctl "brightloop/internal/controller/auth"
	"brightloop/internal/middleware"
)

// Routes registers every /api route. Access is decided here, per group, never inside handlers.
func Routes(s *ghttp.Server) {
	s.Group("/api", func(api *ghttp.RouterGroup) {
		api.Middleware(middleware.Response)
		api.Bind(authctl.Public{})
		api.Group("/", func(in *ghttp.RouterGroup) {
			in.Middleware(middleware.Auth)
			in.Bind(authctl.Session{})
			in.Group("/", func(g *ghttp.RouterGroup) {
				g.Middleware(middleware.RequireRole("admin"))
				g.Bind(adminctl.AdminOnly{})
				if os.Getenv("APP_ENV") == "development" {
					g.Bind(adminctl.DevClock{})
				}
			})
			in.Group("/", func(g *ghttp.RouterGroup) {
				g.Middleware(middleware.RequireRole("admin", "staff"))
				g.Bind(adminctl.BackOffice{})
			})
			in.Group("/", func(g *ghttp.RouterGroup) {
				g.Middleware(middleware.RequireRole("customer"))
				g.Bind(adminctl.Portal{})
			})
		})
	})
}
