package repository

import (
	"context"
	"sync"

	"task-management/internal/domain"
)

// TaskRepository defines the interface for task storage operations
type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	List(ctx context.Context) ([]*domain.Task, error)
	GetTaskByID(ctx context.Context, taskID string) *domain.Task
}

// InMemoryTaskRepository implements TaskRepository interface with in-memory storage
type InMemoryTaskRepository struct {
	tasks map[string]*domain.Task
	mutex sync.RWMutex
}

// NewInMemoryTaskRepository creates a new instance of InMemoryTaskRepository
func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

// Create stores a new task in memory
func (r *InMemoryTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.tasks[task.ID] = task
	return nil
}

// List returns all tasks
func (r *InMemoryTaskRepository) List(ctx context.Context) ([]*domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *InMemoryTaskRepository) GetTaskByID(ctx context.Context, taskID string) *domain.Task {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	task := r.tasks[taskID]
	return task
}
