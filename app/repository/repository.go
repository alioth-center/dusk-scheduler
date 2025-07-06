package repository

import (
	"context"
	"github.com/alioth-center/dusk-scheduler/app/domain"
)

type ClientDao interface {
	CreateClient(ctx context.Context, client *domain.Client) (err error)
}
