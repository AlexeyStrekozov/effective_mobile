package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	core_http_server "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/server"
	"github.com/google/uuid"
)

type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userID *uuid.UUID,
		serviceName *string,
		from *time.Time,
		to *time.Time,
	) (int, error)
}

func NewStatisticsHTTPHandler(
	statisticsService StatisticsService,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: statisticsService,
	}
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
