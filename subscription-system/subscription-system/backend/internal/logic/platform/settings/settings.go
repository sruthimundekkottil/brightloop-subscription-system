// Package settings reads the single company_settings row.
package settings

import (
	"context"
	"errors"
	"time"

	"github.com/gogf/gf/v2/database/gdb"

	"brightloop/internal/dao"
	"brightloop/internal/logic/platform"
	"brightloop/internal/logic/platform/audit"
	"brightloop/internal/model/entity"
)

var ErrMissing = errors.New("company_settings row missing: run the seed")

func Get(ctx context.Context) (*entity.CompanySettings, error) {
	var s *entity.CompanySettings
	if err := dao.CompanySettings.Ctx(ctx).Where("id", 1).Scan(&s); err != nil {
		return nil, err
	}
	if s == nil {
		return nil, ErrMissing
	}
	return s, nil
}

// Update writes the given company_settings columns (nil = NULL) with an audit entry.
func Update(ctx context.Context, set map[string]any) error {
	return platform.Tx(ctx, func(ctx context.Context, tx gdb.TX) error {
		found, err := audit.Update(ctx, tx, "company_settings", 1, set)
		if err == nil && !found {
			err = ErrMissing
		}
		return err
	})
}

// SetTestClock moves the development clock forward and audits it in the same transaction.
func SetTestClock(ctx context.Context, date time.Time) error {
	return platform.Tx(ctx, func(ctx context.Context, tx gdb.TX) error {
		// platform.SetTestClock uses g.DB() with this ctx, which GoFrame routes into tx.
		if err := platform.SetTestClock(ctx, date); err != nil {
			return err
		}
		return audit.Log(ctx, tx, "company_settings", 1, "test-clock", "", map[string]any{"test_clock_date": date.Format(platform.DateLayout)})
	})
}
