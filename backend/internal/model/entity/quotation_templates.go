// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// QuotationTemplates is the golang structure for table quotation_templates.
type QuotationTemplates struct {
	Id             uint64      `json:"id"             ` //
	Name           string      `json:"name"           ` //
	PlanId         uint64      `json:"planId"         ` //
	ValidityDays   uint        `json:"validityDays"   ` //
	PaymentTermsId uint64      `json:"paymentTermsId" ` //
	Terms          string      `json:"terms"          ` //
	Notes          string      `json:"notes"          ` //
	IsActive       int         `json:"isActive"       ` //
	CreatedBy      uint64      `json:"createdBy"      ` //
	CreatedAt      *gtime.Time `json:"createdAt"      ` //
}
