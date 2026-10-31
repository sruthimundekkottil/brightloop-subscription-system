package main

import (
	"os"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"brightloop/internal/cmd"
	"brightloop/internal/cronjob"
	"brightloop/internal/logic/auth"
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
	if auth.Secret = []byte(os.Getenv("JWT_SECRET")); len(auth.Secret) == 0 {
		auth.Secret = []byte(g.Cfg().MustGet(ctx, "jwt.secret").String())
	}
	if len(auth.Secret) < 32 {
		g.Log().Fatal(ctx, "set JWT_SECRET or jwt.secret to at least 32 characters")
	}
	if err := cronjob.Register(ctx); err != nil {
		g.Log().Fatal(ctx, err)
	}
	s := g.Server()
	cmd.Routes(s)
	s.Run()
}
