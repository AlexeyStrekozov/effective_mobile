package subscriptions_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
	core_errors "github.com/AlexeyStrekozov/effective_mobile/internal/core/errors"
	core_postgres_pool "github.com/AlexeyStrekozov/effective_mobile/internal/core/repositorty/postgres/pool"
)

func (r *SubscriptionsRepository) GetSubscription(
	ctx context.Context,
	id int,
) (domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM app.subscriptions
		WHERE id = $1;
	`

	var m SubscriptionModel
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.ServiceName,
		&m.Price,
		&m.UserID,
		&m.StartDate,
		&m.EndDate,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Subscription{}, fmt.Errorf(
				"subscription with id='%d': %w",
				id,
				core_errors.ErrNotFound,
			)
		}
		return domain.Subscription{}, fmt.Errorf("scan get subscription: %w", err)
	}

	return subscriptionDomainFromModel(m), nil
}
