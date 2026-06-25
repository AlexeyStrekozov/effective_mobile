package subscriptions_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
	core_errors "github.com/AlexeyStrekozov/effective_mobile/internal/core/errors"
	core_postgres_pool "github.com/AlexeyStrekozov/effective_mobile/internal/core/repositorty/postgres/pool"
)

func (r *SubscriptionsRepository) UpdateSubscription(
	ctx context.Context,
	id int,
	sub domain.Subscription,
) (domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE app.subscriptions
		SET
			service_name=$1,
			price=$2,
			start_date=$3,
			end_date=$4
		WHERE id=$5
		RETURNING id, service_name, price, user_id, start_date, end_date;
	`

	row := r.pool.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.StartDate,
		sub.EndDate,
		id,
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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Subscription{}, fmt.Errorf(
				"subscription with id='%d': %w",
				id,
				core_errors.ErrNotFound,
			)
		}
		return domain.Subscription{}, fmt.Errorf("scan update subscription: %w", err)
	}

	return subscriptionDomainFromModel(m), nil
}
