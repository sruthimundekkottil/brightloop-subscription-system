// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PasswordResetTokens is the golang structure for table password_reset_tokens.
type PasswordResetTokens struct {
	Id        uint64      `json:"id"        ` //
	UserId    uint64      `json:"userId"    ` //
	TokenHash string      `json:"tokenHash" ` //
	ExpiresAt *gtime.Time `json:"expiresAt" ` //
	UsedAt    *gtime.Time `json:"usedAt"    ` //
	CreatedAt *gtime.Time `json:"createdAt" ` //
}
