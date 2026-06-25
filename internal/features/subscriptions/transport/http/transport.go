package subscriptions_transport_http

import (
	"context"
	"net/http"

	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
	core_http_server "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/server"
	"github.com/google/uuid"
)

type SubscriptionsHTTPHandler struct {
	subscriptionsService SubscriptionsService
}

type SubscriptionsService interface {
	CreateSubscription(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	GetSubscription(ctx context.Context, id int) (domain.Subscription, error)
	GetSubscriptions(ctx context.Context, limit *int, offset *int, userID *uuid.UUID) ([]domain.Subscription, error)
	UpdateSubscription(ctx context.Context, id int, patch domain.SubscriptionPatch) (domain.Subscription, error)
	DeleteSubscription(ctx context.Context, id int) error
}

func NewSubscriptionsHTTPHandler(
	subscriptionsService SubscriptionsService,
) *SubscriptionsHTTPHandler {
	return &SubscriptionsHTTPHandler{
		subscriptionsService: subscriptionsService,
	}
}

func (h *SubscriptionsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/subscriptions",
			Handler: h.CreateSubscription,
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscriptions",
			Handler: h.GetSubscriptions,
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscriptions/{id}",
			Handler: h.GetSubscription,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/subscriptions/{id}",
			Handler: h.UpdateSubscription,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/subscriptions/{id}",
			Handler: h.DeleteSubscription,
		},
	}
}
