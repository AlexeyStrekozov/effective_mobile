package subscriptions_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/AlexeyStrekozov/effective_mobile/internal/core/errors"
)

func (r *SubscriptionsRepository) DeleteSubscription(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM app.subscriptions WHERE id=$1;`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("subscription with id='%d': %w", id, core_errors.ErrNotFound)
	}

	return nil
}
