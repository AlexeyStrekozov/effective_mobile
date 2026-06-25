package subscriptions_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
	core_errors "github.com/AlexeyStrekozov/effective_mobile/internal/core/errors"
	core_logger "github.com/AlexeyStrekozov/effective_mobile/internal/core/logger"
	core_http_request "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/request"
	core_http_response "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/response"
	core_http_types "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/types"
	"github.com/google/uuid"
)

type CreateSubscriptionRequest struct {
	ServiceName string                      `json:"service_name"`
	Price       int                         `json:"price"`
	UserID      uuid.UUID                   `json:"user_id"`
	StartDate   core_http_types.MonthYear  `json:"start_date" swaggertype:"string" example:"2025-07-15"`
	EndDate     *core_http_types.MonthYear `json:"end_date"   swaggertype:"string" example:"2026-07-15"`
}

func (r *CreateSubscriptionRequest) Validate() error {
	if len(r.ServiceName) == 0 {
		return fmt.Errorf("service_name is required: %w", core_errors.ErrInvalidArgument)
	}
	if len(r.ServiceName) > 255 {
		return fmt.Errorf("service_name max length is 255: %w", core_errors.ErrInvalidArgument)
	}
	if r.Price < 0 {
		return fmt.Errorf("price must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if r.UserID == uuid.Nil {
		return fmt.Errorf("user_id is required: %w", core_errors.ErrInvalidArgument)
	}
	if r.StartDate.IsZero() {
		return fmt.Errorf("start_date is required: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

type CreateSubscriptionResponse SubscriptionDTOResponse

// CreateSubscription godoc
// @Summary      Создать подписку
// @Description  Создать новую запись о подписке
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request  body      CreateSubscriptionRequest    true  "Тело запроса"
// @Success      201      {object}  CreateSubscriptionResponse   "Успешно созданная подписка"
// @Failure      400      {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      500      {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /api/v1/subscriptions [post]
func (h *SubscriptionsHTTPHandler) CreateSubscription(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateSubscriptionRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	var endDate *time.Time
	if request.EndDate != nil {
		t := request.EndDate.Time
		endDate = &t
	}

	sub := domain.NewSubscriptionUninitialized(
		request.ServiceName,
		request.Price,
		request.UserID,
		request.StartDate.Time,
		endDate,
	)

	sub, err := h.subscriptionsService.CreateSubscription(ctx, sub)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create subscription")
		return
	}

	responseHandler.JSONResponse(CreateSubscriptionResponse(subscriptionDTOFromDomain(sub)), http.StatusCreated)
}
