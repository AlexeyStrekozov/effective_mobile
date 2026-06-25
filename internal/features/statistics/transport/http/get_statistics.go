package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	core_logger "github.com/AlexeyStrekozov/effective_mobile/internal/core/logger"
	core_http_request "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/request"
	core_http_response "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/response"
	"github.com/google/uuid"
)

type GetStatisticsResponse struct {
	TotalCost int `json:"total_cost"`
}

// GetStatistics godoc
// @Summary      Суммарная стоимость подписок
// @Description  Подсчёт суммарной стоимости подписок за выбранный период с фильтрацией по пользователю и названию сервиса
// @Tags         statistics
// @Produce      json
// @Param        user_id       query     string  false  "Фильтрация по UUID пользователя"
// @Param        service_name  query     string  false  "Фильтрация по названию сервиса"
// @Param        from          query     string  false  "Начало периода (включительно), формат: MM-YYYY"
// @Param        to            query     string  false  "Конец периода (не включительно), формат: MM-YYYY"
// @Success      200           {object}  GetStatisticsResponse            "Суммарная стоимость"
// @Failure      400           {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      500           {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /api/v1/statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, serviceName, from, to, err := getStatisticsQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query params")
		return
	}

	totalCost, err := h.statisticsService.GetStatistics(ctx, userID, serviceName, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return
	}

	responseHandler.JSONResponse(GetStatisticsResponse{TotalCost: totalCost}, http.StatusOK)
}

func getStatisticsQueryParams(r *http.Request) (
	userID *uuid.UUID,
	serviceName *string,
	from *time.Time,
	to *time.Time,
	err error,
) {
	userID, err = core_http_request.GetUUIDQueryParam(r, "user_id")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	serviceNameVal := r.URL.Query().Get("service_name")
	if serviceNameVal != "" {
		serviceName = &serviceNameVal
	}

	from, err = core_http_request.GetMonthYearQueryParam(r, "from")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err = core_http_request.GetMonthYearQueryParam(r, "to")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userID, serviceName, from, to, nil
}
