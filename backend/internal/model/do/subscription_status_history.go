// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionStatusHistory is the golang structure of table subscription_status_history for DAO operations like Where/Data.
type SubscriptionStatusHistory struct {
	g.Meta         `orm:"table:subscription_status_history, do:true"`
	Id             interface{} //
	SubscriptionId interface{} //
	FromStatus     interface{} //
	ToStatus       interface{} //
	ChangedBy      interface{} //
	Note           interface{} //
	ChangedAt      *gtime.Time //
}
