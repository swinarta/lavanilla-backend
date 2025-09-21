package backoffice

import "time"

type DesignerJob struct {
	StartAt time.Time  `json:"start_at"`
	EndAt   *time.Time `json:"end_at,omitempty"`
}
