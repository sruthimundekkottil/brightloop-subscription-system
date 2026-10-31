// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PasswordResetTokens is the golang structure of table password_reset_tokens for DAO operations like Where/Data.
type PasswordResetTokens struct {
	g.Meta    `orm:"table:password_reset_tokens, do:true"`
	Id        interface{} //
	UserId    interface{} //
	TokenHash interface{} //
	ExpiresAt *gtime.Time //
	UsedAt    *gtime.Time //
	CreatedAt *gtime.Time //
}
