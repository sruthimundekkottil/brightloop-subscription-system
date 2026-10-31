// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AuditLog is the golang structure for table audit_log.
type AuditLog struct {
	Id        uint64      `json:"id"        ` //
	TableName string      `json:"tableName" ` //
	RecordId  uint64      `json:"recordId"  ` //
	Action    string      `json:"action"    ` //
	ChangedBy uint64      `json:"changedBy" ` //
	Reason    string      `json:"reason"    ` //
	Changes   string      `json:"changes"   ` //
	CreatedAt *gtime.Time `json:"createdAt" ` //
}
