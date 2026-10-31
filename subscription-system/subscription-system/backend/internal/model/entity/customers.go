// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Customers is the golang structure for table customers.
type Customers struct {
	Id             uint64      `json:"id"             ` //
	UserId         uint64      `json:"userId"         ` //
	Name           string      `json:"name"           ` //
	CompanyName    string      `json:"companyName"    ` //
	Phone          string      `json:"phone"          ` //
	BillingAddress string      `json:"billingAddress" ` //
	CreatedAt      *gtime.Time `json:"createdAt"      ` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      ` //
}
