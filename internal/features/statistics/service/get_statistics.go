package statistics_service

import (
	"context"
	"fmt"
	"time"

	core_errors "github.com/AlexeyStrekozov/effective_mobile/internal/core/errors"
	"github.com/google/uuid"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
	from *time.Time,
	to *time.Time,
) (int, error) {
	if from != nil && to != nil {
		if !to.After(*from) {
			return 0, fmt.Errorf(
				"`to` must be after `from`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	totalCost, err := s.statisticsRepository.GetTotalCost(ctx, userID, serviceName, from, to)
	if err != nil {
		return 0, fmt.Errorf("get total cost from repository: %w", err)
	}

	return totalCost, nil
}
