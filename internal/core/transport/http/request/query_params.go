package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	core_errors "github.com/AlexeyStrekozov/effective_mobile/internal/core/errors"
	"github.com/google/uuid"
)

func GetQueryPrams(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}

func GetDateQueryPram(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	layout := "2006-01-02"
	date, err := time.Parse(layout, param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid date: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &date, nil
}

// GetMonthYearQueryParam parses a query param in "MM-YYYY" format (e.g. "07-2025").
func GetMonthYearQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	t, err := time.Parse("01-2006", param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid MM-YYYY date: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &t, nil
}

func GetUUIDQueryParam(r *http.Request, key string) (*uuid.UUID, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	id, err := uuid.Parse(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid UUID: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &id, nil
}
