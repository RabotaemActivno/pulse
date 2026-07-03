package dto

import (
	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/google/uuid"
)

type UpdateMonitorInput struct {
	ID              uuid.UUID `json:"id"`
	Name            *string   `json:"name"`
	URL             *string   `json:"url"`
	Method          *string   `json:"method"`
	IntervalSeconds *int      `json:"interval_seconds"`
	TimeoutSeconds  *int      `json:"timeout_seconds"`
	ExpectedStatus  *int      `json:"expected_status"`
}

func (i UpdateMonitorInput) Validate() error {
	if i.Name == nil &&
		i.URL == nil &&
		i.Method == nil &&
		i.IntervalSeconds == nil &&
		i.TimeoutSeconds == nil &&
		i.ExpectedStatus == nil {
		return domain.ErrorAllFieldsForUpdateEmpty
	}
	return nil
}
