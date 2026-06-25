package statistics_service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTotalCost(
		ctx context.Context,
		userID *uuid.UUID,
		serviceName *string,
		from *time.Time,
		to *time.Time,
	) (int, error)
}

func NewStatisticsService(
	statisticsRepository StatisticsRepository,
) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
