package subscriptions_service

import (
	"context"
	"fmt"

	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
)

func (s *SubscriptionsService) UpdateSubscription(
	ctx context.Context,
	id int,
	patch domain.SubscriptionPatch,
) (domain.Subscription, error) {
	sub, err := s.subscriptionsRepository.GetSubscription(ctx, id)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("get subscription: %w", err)
	}

	if err := sub.ApplyPatch(patch); err != nil {
		return domain.Subscription{}, fmt.Errorf("apply subscription patch: %w", err)
	}

	updated, err := s.subscriptionsRepository.UpdateSubscription(ctx, id, sub)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("update subscription: %w", err)
	}

	return updated, nil
}
