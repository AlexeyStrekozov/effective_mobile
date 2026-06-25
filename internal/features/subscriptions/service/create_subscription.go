package subscriptions_service

import (
	"context"
	"fmt"

	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
)

func (s *SubscriptionsService) CreateSubscription(
	ctx context.Context,
	sub domain.Subscription,
) (domain.Subscription, error) {
	if err := sub.Validate(); err != nil {
		return domain.Subscription{}, fmt.Errorf("validate subscription: %w", err)
	}

	sub, err := s.subscriptionsRepository.CreateSubscription(ctx, sub)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("create subscription: %w", err)
	}

	return sub, nil
}
