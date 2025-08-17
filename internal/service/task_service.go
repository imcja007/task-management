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
		UpdatedAt:   time.Now(),
	}

	resp, err := s.repo.Create(ctx, task)
	if err != nil {
		return nil, err
	}
	log.Println("Successfully inserted task with ID: ", resp)
	return task, nil
}

// ListTasks returns all tasks
func (s *TaskService) ListTasks(ctx context.Context) ([]*domain.Task, error) {
	return s.repo.List(ctx)
}
func (s *TaskService) GetTaskByIDFromDB(ctx context.Context, taskID string) (*domain.Task, error) {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	return task, nil
}
