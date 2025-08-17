package service

import (
	"context"
	"log"
	"time"

	"task-management/internal/domain"
	"task-management/internal/repository"

	"github.com/google/uuid"
)

// TaskService handles business logic for tasks
type TaskService struct {
	repo repository.TaskRepository
}

// NewTaskService creates a new instance of TaskService
func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(ctx context.Context, title, description string) (*domain.Task, error) {
	task := &domain.Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	resp, err := s.repo.Create(ctx, task)
	if err != nil {
		return nil, err
	}
	log.Println("Successfully inserted task with ID: ", resp)
	return task, nil
}

// ListTasks returns tasks, optionally filtered by status, with pagination
func (s *TaskService) ListTasks(ctx context.Context, status string, page, pageSize int) ([]*domain.Task, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // default page size
	}
	return s.repo.List(ctx, status, page, pageSize)
}

func (s *TaskService) GetTaskByIDFromDB(ctx context.Context, taskID string) (*domain.Task, error) {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID string, updates domain.TaskUpdate) error {
	return s.repo.UpdateInDb(ctx, taskID, updates)
}

// DeleteTask removes a task
func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	return s.repo.DeleteFromDb(ctx, id)
}

func (s *TaskService) UpdateTaskStatus(ctx context.Context, id, status string) error {
	const (
		StatusPending    = "pending"
		StatusInProgress = "in_progress"
		StatusCompleted  = "completed"
	)

	// Validate status
	validStatus := status == StatusPending || status == StatusInProgress || status == StatusCompleted
	if !validStatus {
		return domain.ErrInvalidNewStatus
	}

	// Create update with only status field
	updates := domain.TaskUpdate{
		Status: &status,
	}

	return s.repo.UpdateInDb(ctx, id, updates)
}
