package repository

import (
	"context"
	"errors"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"gorm.io/gorm"
)

type clientDao struct {
	db *gorm.DB
}

func NewClientDao(db *gorm.DB) ClientDao {
	return &clientDao{db: db}
}

func (dao *clientDao) CreateClient(ctx context.Context, client *domain.Client) (clientID uint64, err error) {
	if client == nil {
		return 0, gorm.ErrInvalidData
	}

	if insertErr := dao.db.WithContext(ctx).Model(&domain.Client{}).Create(&client).Error; insertErr != nil {
		return 0, insertErr
	}

	return client.ID, nil
}

func (dao *clientDao) GetClientByID(ctx context.Context, clientID uint64) (client *domain.Client, exist bool, err error) {
	findErr := dao.db.WithContext(ctx).Model(&domain.Client{}).Where("id = ?", clientID).Find(&client).Error
	if findErr == nil {
		return client, true, nil
	}
	if errors.Is(findErr, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	return nil, false, findErr
}
