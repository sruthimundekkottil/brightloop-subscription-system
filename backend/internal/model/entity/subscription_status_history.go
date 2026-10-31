// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionStatusHistory is the golang structure for table subscription_status_history.
type SubscriptionStatusHistory struct {
	Id             uint64      `json:"id"             ` //
	SubscriptionId uint64      `json:"subscriptionId" ` //
	FromStatus     string      `json:"fromStatus"     ` //
	ToStatus       string      `json:"toStatus"       ` //
	ChangedBy      uint64      `json:"changedBy"      ` //
	Note           string      `json:"note"           ` //
	ChangedAt      *gtime.Time `json:"changedAt"      ` //
}
