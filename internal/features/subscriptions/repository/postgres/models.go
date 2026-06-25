package subscriptions_postgres_repository

import (
	"time"

	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
	"github.com/google/uuid"
)

type SubscriptionModel struct {
	ID          int
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

func subscriptionDomainFromModel(m SubscriptionModel) domain.Subscription {
	return domain.NewSubscription(
		m.ID,
		m.ServiceName,
		m.Price,
		m.UserID,
		m.StartDate,
		m.EndDate,
	)
}

func subscriptionDomainFromModels(models []SubscriptionModel) []domain.Subscription {
	result := make([]domain.Subscription, len(models))
	for i, m := range models {
		result[i] = subscriptionDomainFromModel(m)
	}
	return result
}
