// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Payments is the golang structure of table payments for DAO operations like Where/Data.
type Payments struct {
	g.Meta         `orm:"table:payments, do:true"`
	Id             interface{} //
	InvoiceId      interface{} //
	Method         interface{} //
	Amount         interface{} //
	PaymentDate    *gtime.Time //
	Note           interface{} //
	IdempotencyKey interface{} //
	RecordedBy     interface{} //
	CreatedAt      *gtime.Time //
}
