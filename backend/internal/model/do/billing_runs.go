// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// BillingRuns is the golang structure of table billing_runs for DAO operations like Where/Data.
type BillingRuns struct {
	g.Meta       `orm:"table:billing_runs, do:true"`
	Id           interface{} //
	AsOfDate     *gtime.Time //
	TriggeredBy  interface{} //
	DryRun       interface{} //
	StartedAt    *gtime.Time //
	FinishedAt   *gtime.Time //
	DueCount     interface{} //
	CreatedCount interface{} //
	SkippedCount interface{} //
	FailedCount  interface{} //
	Errors       interface{} //
}
