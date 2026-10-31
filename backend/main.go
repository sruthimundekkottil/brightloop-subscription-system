package main

import (
	"os"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"

	"brightloop/internal/cronjob"
	"brightloop/internal/middleware"
	"brightloop/internal/seed"
)

// Usage: `go run .` serves the API and the scheduler; `go run . seed` loads reference data and demo logins.
func main() {
	ctx := gctx.GetInitCtx()
	if link := os.Getenv("DB_LINK"); link != "" { // docker-compose points the API at host "mysql"
		gdb.SetConfig(gdb.Config{"default": gdb.ConfigGroup{{Link: link}}})
	}
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		if err := seed.Run(ctx); err != nil {
			g.Log().Fatal(ctx, err)
		}
		g.Log().Info(ctx, "seed done")
		return
	}
	if err := cronjob.Register(ctx); err != nil {
		g.Log().Fatal(ctx, err)
	}
	s := g.Server()
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.Response)
	})
	s.Run()
}
