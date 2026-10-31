// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// NumberSequences is the golang structure of table number_sequences for DAO operations like Where/Data.
type NumberSequences struct {
	g.Meta    `orm:"table:number_sequences, do:true"`
	Name      interface{} //
	NextValue interface{} //
}
