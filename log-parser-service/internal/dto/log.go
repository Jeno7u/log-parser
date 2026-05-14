package dto

import "time"

type Log struct {
	ID           string    `json:"log_id"`
	FileName     string    `json:"file_name"`
	SourcePath   string    `json:"source_path"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
