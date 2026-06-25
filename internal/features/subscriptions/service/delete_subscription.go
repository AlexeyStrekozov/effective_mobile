package subscriptions_service

import (
	"context"
	"fmt"
)

func (s *SubscriptionsService) DeleteSubscription(
	ctx context.Context,
	id int,
) error {
	if err := s.subscriptionsRepository.DeleteSubscription(ctx, id); err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}

	return nil
}
