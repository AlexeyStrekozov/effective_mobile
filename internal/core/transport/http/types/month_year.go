package core_http_types

import (
	"encoding/json"
	"fmt"
	"time"
)

// MonthYear marshals/unmarshals time.Time as "YYYY-MM-DD" (e.g. "2025-07-15").
type MonthYear struct {
	time.Time
}

func (m *MonthYear) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD (e.g. \"2025-07-15\"): %w", err)
	}

	m.Time = t
	return nil
}

func (m MonthYear) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Time.Format("2006-01-02"))
}
