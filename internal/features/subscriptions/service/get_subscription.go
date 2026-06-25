package subscriptions_service

import (
	"context"
	"fmt"

	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
)

func (s *SubscriptionsService) GetSubscription(
	ctx context.Context,
	id int,
) (domain.Subscription, error) {
	sub, err := s.subscriptionsRepository.GetSubscription(ctx, id)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("get subscription from repository: %w", err)
	}

	return sub, nil
}
