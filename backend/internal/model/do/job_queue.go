// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// JobQueue is the golang structure of table job_queue for DAO operations like Where/Data.
type JobQueue struct {
	g.Meta      `orm:"table:job_queue, do:true"`
	Id          interface{} //
	QueueName   interface{} //
	DedupeKey   interface{} //
	Payload     interface{} //
	EffectiveAt *gtime.Time //
	State       interface{} //
	Attempts    interface{} //
	MaxAttempts interface{} //
	LockedBy    interface{} //
	LockedUntil *gtime.Time //
	LastError   interface{} //
	CreatedAt   *gtime.Time //
	ProcessedAt *gtime.Time //
}
