// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// BillingRuns is the golang structure for table billing_runs.
type BillingRuns struct {
	Id           uint64      `json:"id"           ` //
	AsOfDate     *gtime.Time `json:"asOfDate"     ` //
	TriggeredBy  uint64      `json:"triggeredBy"  ` //
	DryRun       int         `json:"dryRun"       ` //
	StartedAt    *gtime.Time `json:"startedAt"    ` //
	FinishedAt   *gtime.Time `json:"finishedAt"   ` //
	DueCount     uint        `json:"dueCount"     ` //
	CreatedCount uint        `json:"createdCount" ` //
	SkippedCount uint        `json:"skippedCount" ` //
	FailedCount  uint        `json:"failedCount"  ` //
	Errors       string      `json:"errors"       ` //
}
