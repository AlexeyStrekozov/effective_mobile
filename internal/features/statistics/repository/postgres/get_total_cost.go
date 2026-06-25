package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (r *StatisticsRepository) GetTotalCost(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
	from *time.Time,
	to *time.Time,
) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`SELECT COALESCE(SUM(price), 0) FROM app.subscriptions`)

	args := []any{}
	conditions := []string{}

	if userID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id=$%d", len(args)+1))
		args = append(args, userID)
	}

	if serviceName != nil {
		conditions = append(conditions, fmt.Sprintf("service_name=$%d", len(args)+1))
		args = append(args, serviceName)
	}

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("start_date>=$%d", len(args)+1))
		args = append(args, from)
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("start_date<$%d", len(args)+1))
		args = append(args, to)
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	var totalCost int
	if err := r.pool.QueryRow(ctx, queryBuilder.String(), args...).Scan(&totalCost); err != nil {
		return 0, fmt.Errorf("scan total cost: %w", err)
	}

	return totalCost, nil
}
