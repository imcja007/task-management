package repository

import (
	"context"
	"log"
	"os"
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
	List(ctx context.Context, status string, page, pageSize int) ([]*domain.Task, error)
	GetTaskByID(ctx context.Context, taskID string) (*domain.Task, error)
	UpdateInDb(ctx context.Context, taskID string, updates domain.TaskUpdate) error
	DeleteFromDb(ctx context.Context, taskID string) error
}

// InMemoryTaskRepository implements TaskRepository interface with in-memory storage
type InMemoryTaskRepository struct{}

// NewInMemoryTaskRepository creates a new instance of InMemoryTaskRepository
func NewInMemoryTaskRepository() *InMemoryTaskRepository {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	tasksDB = client.Database("taskdb").Collection("tasks")
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

// List returns tasks, optionally filtered by status
func (r *InMemoryTaskRepository) List(ctx context.Context, status string, page, pageSize int) ([]*domain.Task, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	// Calculate skip value for pagination
	skip := (page - 1) * pageSize

	// Set up the options for pagination
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "createdat", Value: -1}}) // Sort by creation date, newest first

	cursor, err := tasksDB.Find(ctx, filter, opts)
	if err != nil {
		log.Println(err)
		return nil, domain.ErrTaskNotFound
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
	if err := tasksDB.FindOne(ctx, bson.D{{Key: "id", Value: taskID}}).Decode(&result); err != nil {
		log.Println(err)
		return nil, domain.ErrTaskNotFound
	}
	return result, nil
}

func (r *InMemoryTaskRepository) DeleteFromDb(ctx context.Context, taskID string) error {
	_, err := tasksDB.DeleteOne(ctx, bson.D{{Key: "id", Value: taskID}})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *InMemoryTaskRepository) UpdateInDb(ctx context.Context, taskID string, updates domain.TaskUpdate) error {
	filter := bson.D{{Key: "id", Value: taskID}}

	// Convert updates to bson.D
	updateFields := bson.D{}

	if updates.Title != nil {
		updateFields = append(updateFields, bson.E{Key: "title", Value: *updates.Title})
	}
	if updates.Description != nil {
		updateFields = append(updateFields, bson.E{Key: "description", Value: *updates.Description})
	}
	if updates.Status != nil {
		updateFields = append(updateFields, bson.E{Key: "status", Value: *updates.Status})
	}

	// update the updated_at timestamp
	updateFields = append(updateFields, bson.E{Key: "updated_at", Value: time.Now()})

	update := bson.D{{Key: "$set", Value: updateFields}}

	result, err := tasksDB.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating task: %v\n", err)
		return err
	}

	if result.MatchedCount == 0 {
		log.Printf("No task found with ID: %s\n", taskID)
		return domain.ErrTaskNotFound
	} else {
		log.Println("Successfully updated task with ID: ", taskID)
	}

	return nil
}
