package subscriptions_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
)

func (r *SubscriptionsRepository) CreateSubscription(
	ctx context.Context,
	sub domain.Subscription,
) (domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO app.subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service_name, price, user_id, start_date, end_date;
	`

	row := r.pool.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)

	var m SubscriptionModel
	err := row.Scan(
		&m.ID,
		&m.ServiceName,
		&m.Price,
		&m.UserID,
		&m.StartDate,
		&m.EndDate,
	)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("scan create subscription: %w", err)
	}

	return subscriptionDomainFromModel(m), nil
}
