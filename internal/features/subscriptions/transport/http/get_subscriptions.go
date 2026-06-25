package subscriptions_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/AlexeyStrekozov/effective_mobile/internal/core/logger"
	core_http_request "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/request"
	core_http_response "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/response"
	"github.com/google/uuid"
)

type GetSubscriptionsResponse []SubscriptionDTOResponse

// GetSubscriptions godoc
// @Summary      Список подписок
// @Description  Список подписок с опциональной пагинацией и/или фильтрацией по user_id
// @Tags         subscriptions
// @Produce      json
// @Param        user_id  query     string  false  "Фильтрация по UUID пользователя"
// @Param        limit    query     int     false  "Размер страницы"
// @Param        offset   query     int     false  "Смещение"
// @Success      200      {object}  GetSubscriptionsResponse         "Список подписок"
// @Failure      400      {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      500      {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /api/v1/subscriptions [get]
func (h *SubscriptionsHTTPHandler) GetSubscriptions(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, limit, offset, err := getSubscriptionsQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query params")
		return
	}

	subs, err := h.subscriptionsService.GetSubscriptions(ctx, limit, offset, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get subscriptions")
		return
	}

	responseHandler.JSONResponse(GetSubscriptionsResponse(subscriptionDTOFromDomains(subs)), http.StatusOK)
}

func getSubscriptionsQueryParams(r *http.Request) (
	_ *uuid.UUID,
	limit *int,
	offset *int,
	_ error,
) {
	userID, err := core_http_request.GetUUIDQueryParam(r, "user_id")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	limit, err = core_http_request.GetQueryPrams(r, "limit")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err = core_http_request.GetQueryPrams(r, "offset")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userID, limit, offset, nil
}
