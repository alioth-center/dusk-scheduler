package repository

import (
	"context"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"gorm.io/gorm"
	"time"
)

type taskDao struct {
	db *gorm.DB
}

func NewTaskDao(db *gorm.DB) TaskDao {
	return &taskDao{db: db}
}

func (dao *taskDao) StatisticsClientQuotaUsage(ctx context.Context, clientID uint64, startTime time.Time) (usage uint64, err error) {
	if queryErr := dao.db.WithContext(ctx).Model(&domain.Task{}).
		Where("submitter = ? and completed_at > ?", clientID, startTime).
		Select("sum(quota_usage)").Scan(&usage).Error; queryErr != nil {
		return 0, queryErr
	}

	return usage, nil
}
