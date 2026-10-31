// Package settings reads the single company_settings row.
package settings

import (
	"context"
	"errors"

	"brightloop/internal/dao"
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
