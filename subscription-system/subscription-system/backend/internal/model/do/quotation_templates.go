// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// QuotationTemplates is the golang structure of table quotation_templates for DAO operations like Where/Data.
type QuotationTemplates struct {
	g.Meta         `orm:"table:quotation_templates, do:true"`
	Id             interface{} //
	Name           interface{} //
	PlanId         interface{} //
	ValidityDays   interface{} //
	PaymentTermsId interface{} //
	Terms          interface{} //
	Notes          interface{} //
	IsActive       interface{} //
	CreatedBy      interface{} //
	CreatedAt      *gtime.Time //
}
