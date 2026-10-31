// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/shopspring/decimal"
)

// TaxRules is the golang structure for table tax_rules.
type TaxRules struct {
	Id            uint64          `json:"id"            ` //
	Name          string          `json:"name"          ` //
	RatePercent   decimal.Decimal `json:"ratePercent"   ` //
	IsDefault     int             `json:"isDefault"     ` //
	IsActive      int             `json:"isActive"      ` //
	DefaultMarker uint            `json:"defaultMarker" ` //
}
