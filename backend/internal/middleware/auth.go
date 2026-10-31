package middleware

import (
	"slices"

	"github.com/gogf/gf/v2/net/ghttp"

	"brightloop/internal/apperr"
	"brightloop/internal/logic/auth"
	"brightloop/internal/logic/platform"
)

// Auth requires a valid session cookie and puts the user into the request context (platform.ActorFrom, auth.CustomerID).
func Auth(r *ghttp.Request) {
	uid, err := auth.ParseToken(r.Cookie.Get(auth.CookieName).String())
	if err == nil {
		var a platform.Actor
		if a, err = auth.LoadActor(r.Context(), uid); err == nil {
			r.SetCtx(platform.WithActor(r.Context(), a))
			r.Middleware.Next()
			return
		}
	}
	r.SetError(err)
}

// RequireRole lets only the given roles through; it must run after Auth.
func RequireRole(roles ...string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		if !slices.Contains(roles, platform.ActorFrom(r.Context()).Role) {
			r.SetError(apperr.Forbidden())
			return
		}
		r.Middleware.Next()
	}
}
