// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CloseReasons is the golang structure of table close_reasons for DAO operations like Where/Data.
type CloseReasons struct {
	g.Meta  `orm:"table:close_reasons, do:true"`
	Id      interface{} //
	Name    interface{} //
	IsChurn interface{} //
}
