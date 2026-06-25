package subscriptions_service

import (
	"context"
	"fmt"

	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
	core_errors "github.com/AlexeyStrekozov/effective_mobile/internal/core/errors"
	"github.com/google/uuid"
)

func (s *SubscriptionsService) GetSubscriptions(
	ctx context.Context,
	limit *int,
	offset *int,
	userID *uuid.UUID,
) ([]domain.Subscription, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	subs, err := s.subscriptionsRepository.GetSubscriptions(ctx, limit, offset, userID)
	if err != nil {
		return nil, fmt.Errorf("get subscriptions: %w", err)
	}

	return subs, nil
}
