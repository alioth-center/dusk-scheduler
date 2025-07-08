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

func (dao *taskDao) CreateTask(ctx context.Context, task *domain.Task) (taskID uint64, err error) {
	//TODO implement me
	panic("implement me")
}

func (dao *taskDao) GetTaskByID(ctx context.Context, taskID uint64) (task *domain.Task, exist bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (dao *taskDao) GetTaskListByClientID(ctx context.Context, clientID uint64, statusFilter []string, offsetTaskID uint64, pageLimit uint32, desc bool) (tasks []*domain.Task, err error) {
	//TODO implement me
	panic("implement me")
}

func (dao *taskDao) UpdateTaskAsCompleted(ctx context.Context, taskID uint64) error {
	//TODO implement me
	panic("implement me")
}

func (dao *taskDao) UpdateTaskAsArchived(ctx context.Context, taskID uint64, reason domain.TaskArchiveReason) error {
	//TODO implement me
	panic("implement me")
}

func (dao *taskDao) StatisticsClientQuotaUsage(ctx context.Context, clientID uint64, startTime time.Time) (usage uint64, err error) {
	if queryErr := dao.db.WithContext(ctx).Model(&domain.Task{}).
		Where("submitter = ? and completed_at > ?", clientID, startTime).
		Select("sum(quota_usage)").Scan(&usage).Error; queryErr != nil {
		return 0, queryErr
	}

	return usage, nil
}
