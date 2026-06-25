package subscriptions_service

import (
	"context"

	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
	"github.com/google/uuid"
)

type SubscriptionsService struct {
	subscriptionsRepository SubscriptionsRepository
}

type SubscriptionsRepository interface {
	CreateSubscription(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	GetSubscription(ctx context.Context, id int) (domain.Subscription, error)
	GetSubscriptions(ctx context.Context, limit *int, offset *int, userID *uuid.UUID) ([]domain.Subscription, error)
	UpdateSubscription(ctx context.Context, id int, sub domain.Subscription) (domain.Subscription, error)
	DeleteSubscription(ctx context.Context, id int) error
}

func NewSubscriptionsService(
	subscriptionsRepository SubscriptionsRepository,
) *SubscriptionsService {
	return &SubscriptionsService{
		subscriptionsRepository: subscriptionsRepository,
	}
}
