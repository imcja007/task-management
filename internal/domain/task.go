package domain

import (
	"time"
)

// Task represents the task entity
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const (
	StatusPending   = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted = "completed"
)
