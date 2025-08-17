package domain

import "errors"

// Error definitions for the domain package
var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrInvalidNewStatus = errors.New("invalid status. Must be one of pending, in_progress, or completed")
)
