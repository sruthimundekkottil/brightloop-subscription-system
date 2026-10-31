// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/shopspring/decimal"
)

// ProductAttributeValues is the golang structure for table product_attribute_values.
type ProductAttributeValues struct {
	Id               uint64          `json:"id"               ` //
	ProductId        uint64          `json:"productId"        ` //
	AttributeValueId uint64          `json:"attributeValueId" ` //
	PriceExtra       decimal.Decimal `json:"priceExtra"       ` //
	IsActive         int             `json:"isActive"         ` //
}
