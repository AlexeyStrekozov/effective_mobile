package subscriptions_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AlexeyStrekozov/effective_mobile/internal/core/domain"
	"github.com/google/uuid"
)

func (r *SubscriptionsRepository) GetSubscriptions(
	ctx context.Context,
	limit *int,
	offset *int,
	userID *uuid.UUID,
) ([]domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM app.subscriptions
		%s
		ORDER BY id ASC
		LIMIT $1
		OFFSET $2;
	`

	args := []any{limit, offset}

	if userID != nil {
		query = fmt.Sprintf(query, "WHERE user_id=$3")
		args = append(args, userID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select subscriptions: %w", err)
	}
	defer rows.Close()

	var models []SubscriptionModel
	for rows.Next() {
		var m SubscriptionModel
		err := rows.Scan(
			&m.ID,
			&m.ServiceName,
			&m.Price,
			&m.UserID,
			&m.StartDate,
			&m.EndDate,
		)
		if err != nil {
			return nil, fmt.Errorf("scan subscription: %w", err)
		}
		models = append(models, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return subscriptionDomainFromModels(models), nil
}
