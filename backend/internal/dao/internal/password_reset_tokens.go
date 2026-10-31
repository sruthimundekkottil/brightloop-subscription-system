// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PasswordResetTokensDao is the data access object for table password_reset_tokens.
type PasswordResetTokensDao struct {
	table   string                     // table is the underlying table name of the DAO.
	group   string                     // group is the database configuration group name of current DAO.
	columns PasswordResetTokensColumns // columns contains all the column names of Table for convenient usage.
}

// PasswordResetTokensColumns defines and stores column names for table password_reset_tokens.
type PasswordResetTokensColumns struct {
	Id        string //
	UserId    string //
	TokenHash string //
	ExpiresAt string //
	UsedAt    string //
	CreatedAt string //
}

// passwordResetTokensColumns holds the columns for table password_reset_tokens.
var passwordResetTokensColumns = PasswordResetTokensColumns{
	Id:        "id",
	UserId:    "user_id",
	TokenHash: "token_hash",
	ExpiresAt: "expires_at",
	UsedAt:    "used_at",
	CreatedAt: "created_at",
}

// NewPasswordResetTokensDao creates and returns a new DAO object for table data access.
func NewPasswordResetTokensDao() *PasswordResetTokensDao {
	return &PasswordResetTokensDao{
		group:   "default",
		table:   "password_reset_tokens",
		columns: passwordResetTokensColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PasswordResetTokensDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PasswordResetTokensDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PasswordResetTokensDao) Columns() PasswordResetTokensColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PasswordResetTokensDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PasswordResetTokensDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PasswordResetTokensDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
