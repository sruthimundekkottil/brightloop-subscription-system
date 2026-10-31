package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"brightloop/internal/cronjob"
)

// shortcut: runs only the scheduler; the HTTP server joins here when the API is built.
func main() {
	ctx := gctx.GetInitCtx()
	if err := cronjob.Register(ctx); err != nil {
		g.Log().Fatal(ctx, err)
	}
	g.Log().Info(ctx, "scheduler started")
	select {}
}
