// Package auth is the request and response shapes of /api/auth.
package auth

import "github.com/gogf/gf/v2/frame/g"

type RegisterReq struct {
	g.Meta      `path:"/auth/register" method:"post" tags:"Auth" summary:"Register a customer account"`
	Email       string `json:"email"       v:"required|email|max-length:254"`
	Password    string `json:"password"    v:"required"` // the full rule is checked in logic/auth
	Name        string `json:"name"        v:"required|length:1,150"`
	CompanyName string `json:"companyName" v:"max-length:150"`
	Phone       string `json:"phone"       v:"max-length:30"`
}

type LoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" tags:"Auth" summary:"Log in"`
	Email    string `json:"email"    v:"required"`
	Password string `json:"password" v:"required"`
}

type LogoutReq struct {
	g.Meta `path:"/auth/logout" method:"post" tags:"Auth" summary:"Log out"`
}
type LogoutRes struct{}

type MeReq struct {
	g.Meta `path:"/auth/me" method:"get" tags:"Auth" summary:"The logged-in user"`
}

// MeRes is returned by register, login and me.
type MeRes struct {
	ID         uint64  `json:"id"`
	Email      string  `json:"email"`
	Role       string  `json:"role"`
	CustomerID *uint64 `json:"customerId"` // null for staff and admins
}

type ForgotPasswordReq struct {
	g.Meta `path:"/auth/forgot-password" method:"post" tags:"Auth" summary:"Email a reset link (always 200)"`
	Email  string `json:"email" v:"required|email|max-length:254"`
}

type ResetPasswordReq struct {
	g.Meta   `path:"/auth/reset-password" method:"post" tags:"Auth" summary:"Set a new password with a reset token"`
	Token    string `json:"token"    v:"required|max-length:100"`
	Password string `json:"password" v:"required"`
}

// MessageRes is a reply with only a human-readable message.
type MessageRes struct {
	Message string `json:"message"`
}
