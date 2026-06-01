package dto

import "github.com/RabotaemActivno/pulse/internal/domain"

type GetMonitorsOutput struct {
	Monitors []domain.Monitor `json:"monitors"` 
}
