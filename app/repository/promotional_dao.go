package repository

import (
	"context"
	"errors"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"gorm.io/gorm"
)

type promotionalDao struct {
	db *gorm.DB
}

func NewPromotionalDao(db *gorm.DB) PromotionalDao {
	return &promotionalDao{db: db}
}

func (dao *promotionalDao) GetPromotionalByCode(ctx context.Context, code string) (promotional *domain.Promotional, exist bool, err error) {
	condition, result := &domain.Promotional{Code: code}, &domain.Promotional{}
	queryErr := dao.db.WithContext(ctx).Model(&domain.Promotional{}).Where(condition).Find(&result).Error
	if queryErr == nil {
		return result, true, nil
	}
	if errors.Is(queryErr, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	return nil, false, queryErr
}
