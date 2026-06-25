package subscriptions_transport_http

import (
	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
	core_http_types "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/types"
	"github.com/google/uuid"
)

type SubscriptionDTOResponse struct {
	ID          int                         `json:"id"`
	ServiceName string                      `json:"service_name"`
	Price       int                         `json:"price"`
	UserID      uuid.UUID                   `json:"user_id"`
	StartDate   core_http_types.MonthYear  `json:"start_date" swaggertype:"string" example:"2025-07-15"`
	EndDate     *core_http_types.MonthYear `json:"end_date"   swaggertype:"string" example:"2026-07-15"`
}

func subscriptionDTOFromDomain(sub domain.Subscription) SubscriptionDTOResponse {
	dto := SubscriptionDTOResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   core_http_types.MonthYear{Time: sub.StartDate},
	}
	if sub.EndDate != nil {
		my := core_http_types.MonthYear{Time: *sub.EndDate}
		dto.EndDate = &my
	}
	return dto
}

func subscriptionDTOFromDomains(subs []domain.Subscription) []SubscriptionDTOResponse {
	dtos := make([]SubscriptionDTOResponse, len(subs))
	for i, sub := range subs {
		dtos[i] = subscriptionDTOFromDomain(sub)
	}
	return dtos
}
