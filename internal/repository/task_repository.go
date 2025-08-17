package repository

import (
	"context"
	"log"
	"task-management/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tasksDB *mongo.Collection

// TaskRepository defines the interface for task storage operations
type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) (string, error)
	List(ctx context.Context) ([]*domain.Task, error)
	GetTaskByID(ctx context.Context, taskID string) (*domain.Task, error)
}

// InMemoryTaskRepository implements TaskRepository interface with in-memory storage
type InMemoryTaskRepository struct {}

// NewInMemoryTaskRepository creates a new instance of InMemoryTaskRepository
func NewInMemoryTaskRepository() *InMemoryTaskRepository {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(ctx, clientOptions)
	tasksDB = client.Database("task_db").Collection("tasks")
	return &InMemoryTaskRepository{}
}

// Create stores a new task in memory
func (r *InMemoryTaskRepository) Create(ctx context.Context, task *domain.Task) (string, error) {
	_, err := tasksDB.InsertOne(ctx, task)
	if err != nil {
		log.Println("Something went wrong while inserting tasks in db")
		return "", err
	}
	return task.ID, nil
}

// List returns all tasks
func (r *InMemoryTaskRepository) List(ctx context.Context) ([]*domain.Task, error) {
	cursor, err := tasksDB.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Oops something went wrong")
		return nil, err
	}
	var tasks []*domain.Task
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var task domain.Task
		if err := cursor.Decode(&task); err != nil {
			log.Println("Error decoding task:", err)
			continue
		}
		tasks = append(tasks, &task)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return tasks, nil
}

func (r *InMemoryTaskRepository) GetTaskByID(ctx context.Context, taskID string) (*domain.Task, error) {
	var result *domain.Task
	err := tasksDB.FindOne(ctx, bson.D{{Key: "id", Value: taskID}}).Decode(&result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}
