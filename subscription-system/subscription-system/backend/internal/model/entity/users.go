// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure for table users.
type Users struct {
	Id           uint64      `json:"id"           ` //
	Email        string      `json:"email"        ` //
	PasswordHash string      `json:"passwordHash" ` //
	Role         string      `json:"role"         ` //
	IsActive     int         `json:"isActive"     ` //
	CreatedAt    *gtime.Time `json:"createdAt"    ` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    ` //
}
