package subscriptions_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
	core_logger "github.com/AlexeyStrekozov/effective_mobile/internal/core/logger"
	core_http_request "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/request"
	core_http_response "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/response"
	core_http_types "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/types"
)

type UpdateSubscriptionRequest struct {
	ServiceName core_http_types.Nullable[string]                    `json:"service_name"`
	Price       core_http_types.Nullable[int]                       `json:"price"`
	StartDate   core_http_types.Nullable[core_http_types.MonthYear] `json:"start_date" swaggertype:"string" example:"2025-07-15"`
	EndDate     core_http_types.Nullable[core_http_types.MonthYear] `json:"end_date"   swaggertype:"string" example:"2026-07-15"`
}

func (req *UpdateSubscriptionRequest) Validate() error {
	if req.ServiceName.Set {
		if req.ServiceName.Value == nil {
			return fmt.Errorf("service_name cannot be null")
		}
		if len(*req.ServiceName.Value) == 0 || len(*req.ServiceName.Value) > 255 {
			return fmt.Errorf("service_name must be between 1 and 255 characters")
		}
	}
	if req.Price.Set {
		if req.Price.Value == nil {
			return fmt.Errorf("price cannot be null")
		}
		if *req.Price.Value < 0 {
			return fmt.Errorf("price must be non-negative")
		}
	}
	if req.StartDate.Set && req.StartDate.Value == nil {
		return fmt.Errorf("start_date cannot be null")
	}
	return nil
}

type UpdateSubscriptionResponse SubscriptionDTOResponse

// UpdateSubscription godoc
// @Summary      Обновить подписку
// @Description  Частичное обновление записи о подписке (PATCH). Поля не переданные в теле — не изменяются. end_date можно сбросить, передав null.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id       path      int                        true  "ID подписки" Format(int)
// @Param        request  body      UpdateSubscriptionRequest  true  "Тело запроса"
// @Success      200      {object}  UpdateSubscriptionResponse         "Обновлённая подписка"
// @Failure      400      {object}  core_http_response.ErrorResponse   "Bad request"
// @Failure      404      {object}  core_http_response.ErrorResponse   "Not found"
// @Failure      500      {object}  core_http_response.ErrorResponse   "Internal server error"
// @Router       /api/v1/subscriptions/{id} [patch]
func (h *SubscriptionsHTTPHandler) UpdateSubscription(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get subscription id path value")
		return
	}

	var request UpdateSubscriptionRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	patch := subscriptionPatchFromRequest(request)

	sub, err := h.subscriptionsService.UpdateSubscription(ctx, id, patch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to update subscription")
		return
	}

	responseHandler.JSONResponse(UpdateSubscriptionResponse(subscriptionDTOFromDomain(sub)), http.StatusOK)
}

func subscriptionPatchFromRequest(req UpdateSubscriptionRequest) domain.SubscriptionPatch {
	patch := domain.SubscriptionPatch{
		ServiceName: req.ServiceName.ToDomain(),
		Price:       req.Price.ToDomain(),
	}

	patch.StartDate = monthYearNullableToDomain(req.StartDate)
	patch.EndDate = monthYearNullableToDomain(req.EndDate)

	return patch
}

func monthYearNullableToDomain(n core_http_types.Nullable[core_http_types.MonthYear]) domain.Nullable[time.Time] {
	if !n.Set {
		return domain.Nullable[time.Time]{}
	}
	if n.Value == nil {
		return domain.Nullable[time.Time]{Set: true, Value: nil}
	}
	t := n.Value.Time
	return domain.Nullable[time.Time]{Set: true, Value: &t}
}
