package dto

import "github.com/google/uuid"

type CreateMonitorInput struct {
	UserID          uuid.UUID `json:"user_id"`
	Name            string    `json:"name"`
	URL             string    `json:"url"`
	Method          string    `json:"method"`
	IntervalSeconds int       `json:"interval_seconds"`
}

type CreateMonitorOutput struct {
	ID uuid.UUID `json:"id"`
}
