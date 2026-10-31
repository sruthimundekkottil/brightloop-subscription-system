// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// StaffProfiles is the golang structure for table staff_profiles.
type StaffProfiles struct {
	Id        uint64      `json:"id"        ` //
	UserId    uint64      `json:"userId"    ` //
	FullName  string      `json:"fullName"  ` //
	JobTitle  string      `json:"jobTitle"  ` //
	CreatedBy uint64      `json:"createdBy" ` //
	CreatedAt *gtime.Time `json:"createdAt" ` //
}
