// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// StaffProfiles is the golang structure of table staff_profiles for DAO operations like Where/Data.
type StaffProfiles struct {
	g.Meta    `orm:"table:staff_profiles, do:true"`
	Id        interface{} //
	UserId    interface{} //
	FullName  interface{} //
	JobTitle  interface{} //
	CreatedBy interface{} //
	CreatedAt *gtime.Time //
}
