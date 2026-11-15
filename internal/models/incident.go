package models

import "time"

type Incident struct {
	ID             int       `json:"id"`
	DeviceID       int       `json:"device_id"`
	DeviceName     string    `json:"device_name"`
	IP             string    `json:"ip"`
	PreviousStatus string    `json:"previous_status"`
	CurrentStatus  string    `json:"current_status"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
}
