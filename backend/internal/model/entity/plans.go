// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

// Plans is the golang structure for table plans.
type Plans struct {
	Id                   uint64          `json:"id"                   ` //
	Name                 string          `json:"name"                 ` //
	Price                decimal.Decimal `json:"price"                ` //
	BillingPeriod        string          `json:"billingPeriod"        ` //
	BillingIntervalCount uint            `json:"billingIntervalCount" ` //
	MinQuantity          uint            `json:"minQuantity"          ` //
	AvailableFrom        *gtime.Time     `json:"availableFrom"        ` //
	AvailableUntil       *gtime.Time     `json:"availableUntil"       ` //
	AutoClose            int             `json:"autoClose"            ` //
	CanCloseEarly        int             `json:"canCloseEarly"        ` //
	CanPause             int             `json:"canPause"             ` //
	CanRenew             int             `json:"canRenew"             ` //
	IsActive             int             `json:"isActive"             ` //
	CreatedAt            *gtime.Time     `json:"createdAt"            ` //
	UpdatedAt            *gtime.Time     `json:"updatedAt"            ` //
}
