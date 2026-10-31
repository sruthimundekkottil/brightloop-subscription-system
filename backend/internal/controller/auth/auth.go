// Package auth is the thin HTTP layer over logic/auth: it sets and clears the session cookie.
package auth

import (
	"context"
	"net/http"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	api "brightloop/api/auth"
	"brightloop/internal/logic/auth"
	"brightloop/internal/logic/platform"
)

// Public holds the routes that need no session.
type Public struct{}

// Session holds the routes that run behind middleware.Auth.
type Session struct{}

func (Public) Register(ctx context.Context, req *api.RegisterReq) (*api.MeRes, error) {
	a, err := auth.Register(ctx, auth.RegisterInput{Email: req.Email, Password: req.Password, Name: req.Name,
		CompanyName: req.CompanyName, Phone: req.Phone})
	if err != nil {
		return nil, err
	}
	if err = startSession(ctx, a.UserID); err != nil {
		return nil, err
	}
	ghttp.RequestFromCtx(ctx).Response.WriteHeader(http.StatusCreated)
	return me(ctx, a)
}

func (Public) Login(ctx context.Context, req *api.LoginReq) (*api.MeRes, error) {
	a, err := auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	if err = startSession(ctx, a.UserID); err != nil {
		return nil, err
	}
	return me(ctx, a)
}

func (Session) Logout(ctx context.Context, _ *api.LogoutReq) (*api.LogoutRes, error) {
	setCookie(ctx, "", -1)
	return &api.LogoutRes{}, nil
}

func (Session) Me(ctx context.Context, _ *api.MeReq) (*api.MeRes, error) {
	return me(ctx, platform.ActorFrom(ctx))
}

func me(ctx context.Context, a platform.Actor) (*api.MeRes, error) {
	email, err := g.DB().GetValue(ctx, "SELECT email FROM users WHERE id = ?", a.UserID)
	if err != nil {
		return nil, err
	}
	res := &api.MeRes{ID: a.UserID, Email: email.String(), Role: a.Role}
	if a.CustomerID != 0 {
		res.CustomerID = &a.CustomerID
	}
	return res, nil
}

func startSession(ctx context.Context, userID uint64) error {
	tok, err := auth.IssueToken(userID)
	if err != nil {
		return err
	}
	setCookie(ctx, tok, int(auth.SessionTTL.Seconds()))
	return nil
}

// setCookie writes the session cookie: httpOnly so scripts can't read it, SameSite=Lax against cross-site posts,
// Secure everywhere except local development over plain http.
func setCookie(ctx context.Context, value string, maxAge int) {
	ghttp.RequestFromCtx(ctx).Cookie.SetHttpCookie(&http.Cookie{
		Name: auth.CookieName, Value: value, Path: "/", MaxAge: maxAge,
		HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: os.Getenv("APP_ENV") != "development",
	})
}
