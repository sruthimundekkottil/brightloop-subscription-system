// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// StaffProfilesDao is the data access object for table staff_profiles.
type StaffProfilesDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns StaffProfilesColumns // columns contains all the column names of Table for convenient usage.
}

// StaffProfilesColumns defines and stores column names for table staff_profiles.
type StaffProfilesColumns struct {
	Id        string //
	UserId    string //
	FullName  string //
	JobTitle  string //
	CreatedBy string //
	CreatedAt string //
}

// staffProfilesColumns holds the columns for table staff_profiles.
var staffProfilesColumns = StaffProfilesColumns{
	Id:        "id",
	UserId:    "user_id",
	FullName:  "full_name",
	JobTitle:  "job_title",
	CreatedBy: "created_by",
	CreatedAt: "created_at",
}

// NewStaffProfilesDao creates and returns a new DAO object for table data access.
func NewStaffProfilesDao() *StaffProfilesDao {
	return &StaffProfilesDao{
		group:   "default",
		table:   "staff_profiles",
		columns: staffProfilesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *StaffProfilesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *StaffProfilesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *StaffProfilesDao) Columns() StaffProfilesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *StaffProfilesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *StaffProfilesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *StaffProfilesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
