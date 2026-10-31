// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// JobQueueDao is the data access object for table job_queue.
type JobQueueDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns JobQueueColumns // columns contains all the column names of Table for convenient usage.
}

// JobQueueColumns defines and stores column names for table job_queue.
type JobQueueColumns struct {
	Id          string //
	QueueName   string //
	DedupeKey   string //
	Payload     string //
	EffectiveAt string //
	State       string //
	Attempts    string //
	MaxAttempts string //
	LockedBy    string //
	LockedUntil string //
	LastError   string //
	CreatedAt   string //
	ProcessedAt string //
}

// jobQueueColumns holds the columns for table job_queue.
var jobQueueColumns = JobQueueColumns{
	Id:          "id",
	QueueName:   "queue_name",
	DedupeKey:   "dedupe_key",
	Payload:     "payload",
	EffectiveAt: "effective_at",
	State:       "state",
	Attempts:    "attempts",
	MaxAttempts: "max_attempts",
	LockedBy:    "locked_by",
	LockedUntil: "locked_until",
	LastError:   "last_error",
	CreatedAt:   "created_at",
	ProcessedAt: "processed_at",
}

// NewJobQueueDao creates and returns a new DAO object for table data access.
func NewJobQueueDao() *JobQueueDao {
	return &JobQueueDao{
		group:   "default",
		table:   "job_queue",
		columns: jobQueueColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *JobQueueDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *JobQueueDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *JobQueueDao) Columns() JobQueueColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *JobQueueDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *JobQueueDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *JobQueueDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
