package domain

import (
	"fmt"
	"time"

	core_errors "github.com/AlexeyStrekozov/effective_mobile/internal/core/errors"
	"github.com/google/uuid"
)

type Subscription struct {
	ID          int
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

func NewSubscription(
	id int,
	serviceName string,
	price int,
	userID uuid.UUID,
	startDate time.Time,
	endDate *time.Time,
) Subscription {
	return Subscription{
		ID:          id,
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}
}

func NewSubscriptionUninitialized(
	serviceName string,
	price int,
	userID uuid.UUID,
	startDate time.Time,
	endDate *time.Time,
) Subscription {
	return Subscription{
		ID:          UninitializedID,
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}
}

func (s *Subscription) Validate() error {
	if len(s.ServiceName) == 0 {
		return fmt.Errorf("service_name is required: %w", core_errors.ErrInvalidArgument)
	}
	if len(s.ServiceName) > 255 {
		return fmt.Errorf("service_name max length is 255: %w", core_errors.ErrInvalidArgument)
	}
	if s.Price < 0 {
		return fmt.Errorf("price must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if s.EndDate != nil && s.EndDate.Before(s.StartDate) {
		return fmt.Errorf("end_date cannot be before start_date: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

type SubscriptionPatch struct {
	ServiceName Nullable[string]
	Price       Nullable[int]
	StartDate   Nullable[time.Time]
	EndDate     Nullable[time.Time]
}

func (s *Subscription) ApplyPatch(patch SubscriptionPatch) error {
	tmp := *s

	if patch.ServiceName.Set && patch.ServiceName.Value != nil {
		tmp.ServiceName = *patch.ServiceName.Value
	}

	if patch.Price.Set && patch.Price.Value != nil {
		tmp.Price = *patch.Price.Value
	}

	if patch.StartDate.Set && patch.StartDate.Value != nil {
		tmp.StartDate = *patch.StartDate.Value
	}

	if patch.EndDate.Set {
		tmp.EndDate = patch.EndDate.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched subscription: %w", err)
	}

	*s = tmp
	return nil
}
