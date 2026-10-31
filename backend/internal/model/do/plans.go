// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Plans is the golang structure of table plans for DAO operations like Where/Data.
type Plans struct {
	g.Meta               `orm:"table:plans, do:true"`
	Id                   interface{} //
	Name                 interface{} //
	Price                interface{} //
	BillingPeriod        interface{} //
	BillingIntervalCount interface{} //
	MinQuantity          interface{} //
	AvailableFrom        *gtime.Time //
	AvailableUntil       *gtime.Time //
	AutoClose            interface{} //
	CanCloseEarly        interface{} //
	CanPause             interface{} //
	CanRenew             interface{} //
	IsActive             interface{} //
	CreatedAt            *gtime.Time //
	UpdatedAt            *gtime.Time //
}
