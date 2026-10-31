// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// JobQueue is the golang structure for table job_queue.
type JobQueue struct {
	Id          uint64      `json:"id"          ` //
	QueueName   string      `json:"queueName"   ` //
	DedupeKey   string      `json:"dedupeKey"   ` //
	Payload     string      `json:"payload"     ` //
	EffectiveAt *gtime.Time `json:"effectiveAt" ` //
	State       string      `json:"state"       ` //
	Attempts    uint        `json:"attempts"    ` //
	MaxAttempts uint        `json:"maxAttempts" ` //
	LockedBy    string      `json:"lockedBy"    ` //
	LockedUntil *gtime.Time `json:"lockedUntil" ` //
	LastError   string      `json:"lastError"   ` //
	CreatedAt   *gtime.Time `json:"createdAt"   ` //
	ProcessedAt *gtime.Time `json:"processedAt" ` //
}
