package subscriptions_transport_http

import (
	"net/http"

	core_logger "github.com/AlexeyStrekozov/effective_mobile/internal/core/logger"
	core_http_request "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/request"
	core_http_response "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/response"
)

// DeleteSubscription godoc
// @Summary      Удалить подписку
// @Description  Удалить запись о подписке по её ID
// @Tags         subscriptions
// @Param        id   path  int  true "ID подписки" Format(int)
// @Success      204  "Успешное удаление"
// @Failure      400  {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      404  {object}  core_http_response.ErrorResponse "Not found"
// @Failure      500  {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /api/v1/subscriptions/{id} [delete]
func (h *SubscriptionsHTTPHandler) DeleteSubscription(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get subscription id path value")
		return
	}

	if err := h.subscriptionsService.DeleteSubscription(ctx, id); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete subscription")
		return
	}

	responseHandler.NoContentResponse()
}
