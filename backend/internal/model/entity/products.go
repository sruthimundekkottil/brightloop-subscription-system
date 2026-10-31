// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

// Products is the golang structure for table products.
type Products struct {
	Id          uint64          `json:"id"          ` //
	Name        string          `json:"name"        ` //
	ProductType string          `json:"productType" ` //
	SalesPrice  decimal.Decimal `json:"salesPrice"  ` //
	CostPrice   decimal.Decimal `json:"costPrice"   ` //
	Notes       string          `json:"notes"       ` //
	IsRecurring int             `json:"isRecurring" ` //
	IsActive    int             `json:"isActive"    ` //
	CreatedAt   *gtime.Time     `json:"createdAt"   ` //
	UpdatedAt   *gtime.Time     `json:"updatedAt"   ` //
}
