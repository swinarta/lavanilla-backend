package metadata

import "time"

type DesignerJob struct {
	StartAt time.Time  `json:"start_at"`
	EndAt   *time.Time `json:"end_at,omitempty"`
}

type Timeline struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
}
