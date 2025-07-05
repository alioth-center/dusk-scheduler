package service

import (
	"context"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"github.com/alioth-center/dusk-scheduler/infra/apis/location"
	"time"
)

type EmailService interface {
	ValidateEmailAddress(ctx context.Context, emailAddress string) (err error)
	SendEmail(ctx context.Context, receiver string, templateKey string, args map[string]any) (err error)
}

type LocationService interface {
	DetectIPLocation(ctx context.Context, ip string) (address *location.Address, err error)
}

type ClientService interface {
	CreateClient(ctx context.Context, email string, promotionCode string, ip string) (client *domain.Client, err error)
	StoreAuthorizationCode(ctx context.Context, clientID uint64, authorizationCode string) (expiredAt time.Time, err error)
}
