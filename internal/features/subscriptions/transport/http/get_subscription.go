package subscriptions_transport_http

import (
	"net/http"

	core_logger "github.com/AlexeyStrekozov/effective_mobile/internal/core/logger"
	core_http_request "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/request"
	core_http_response "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/response"
)

type GetSubscriptionResponse SubscriptionDTOResponse

// GetSubscription godoc
// @Summary      Получить подписку
// @Description  Получение конкретной подписки по её ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true "ID подписки" Format(int)
// @Success      200  {object}  GetSubscriptionResponse          "Подписка найдена"
// @Failure      400  {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      404  {object}  core_http_response.ErrorResponse "Not found"
// @Failure      500  {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /api/v1/subscriptions/{id} [get]
func (h *SubscriptionsHTTPHandler) GetSubscription(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get subscription id path value")
		return
	}

	sub, err := h.subscriptionsService.GetSubscription(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get subscription")
		return
	}

	responseHandler.JSONResponse(GetSubscriptionResponse(subscriptionDTOFromDomain(sub)), http.StatusOK)
}
