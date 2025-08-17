 # ✅ Task Management System

A RESTful microservice for managing tasks built with Go, MongoDB, and following microservices architecture principles.

## 📑 Table of Contents
- [🔍 Problem Breakdown](#-problem-breakdown)
- [🧩 Design Decisions](#-design-decisions)
- [🏗️ Architecture Overview](#️-architecture-overview)
- [🧠 Microservices Concepts](#-microservices-concepts)
- [🚀 Getting Started](#-getting-started)
- [📘 API Documentation](#-api-documentation)
- [📈 Scalability](#-scalability)
  
## 🔍 Problem Breakdown

The Task Management System addresses the need for a scalable, maintainable service to handle basic CRUD operations on tasks. The key requirements addressed are:

1. **Core Functionality** 🛠️: Create, Read, Update, and Delete tasks
2. **Status Management** 🔄: Track task status (pending, in_progress, completed)
3. **Filtering** 🔍: Filter tasks by status
4. **Pagination** 📄: Handle large datasets efficiently
5. **Microservice Architecture** 🏗️: Implement scalable, maintainable service design
6. **RESTful API** 🌐: Provide consistent, intuitive endpoints

## 🧩 Design Decisions

### 🏛️ Architecture Pattern
- **Clean Architecture**: Separated into distinct layers (domain, service, repository, api)
- **Dependency Injection**: Service and repository dependencies are injected
- **Interface-Based Design**: Repository interface allows for easy testing and implementation swapping

### 💻 Technology Stack
- **Go** 🚀: High-performance, concurrent language ideal for microservices
- **MongoDB** 🗃️: NoSQL database for flexible document storage
- **Gorilla Mux** 🔄: HTTP router for clean route handling
- **UUID** 🔑: Unique identifier generation for tasks

### 📊 Data Model
```go
type Task struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 🔄 Status Management
Tasks support three statuses:
- `pending` ⏳: Default status for new tasks
- `in_progress` 🔄: Tasks currently being worked on
- `completed` ✅: Finished tasks

## 🏗️ Architecture Overview

```
├── cmd/
│   └── server/          # Application entry point
│       └── main.go
├── api/                 # HTTP handlers and routing
│   └── task_handler.go
├── internal/
│   ├── domain/          # Business entities and errors
│   │   ├── task.go
│   │   ├── task_update.go
│   │   └── errors.go
│   ├── service/         # Business logic layer
│   │   └── task_service.go
│   └── repository/      # Data access layer
│       └── task_repository.go
├── go.mod
└── go.sum
```

### 🔍 Layer Responsibilities

1. **API Layer** (`api/`): Handles HTTP requests, response formatting, and routing
2. **Service Layer** (`internal/service/`): Contains business logic and orchestrates operations
3. **Repository Layer** (`internal/repository/`): Manages data persistence and retrieval
4. **Domain Layer** (`internal/domain/`): Defines business entities and domain-specific errors

## 🧠 Microservices Concepts

### 🎯 Single Responsibility Principle
Each component has a single, well-defined responsibility:
- **TaskHandler**: HTTP request/response handling
- **TaskService**: Business logic execution
- **TaskRepository**: Data persistence operations
- **Task Domain**: Entity definition and validation

### 🧩 Separation of Concerns
- HTTP concerns are isolated in the API layer
- Business logic is contained in the service layer
- Data access logic is abstracted in the repository layer
- Domain entities are pure business objects

### 🔌 Interface-Based Design
```go
type TaskRepository interface {
    Create(ctx context.Context, task *domain.Task) (string, error)
    List(ctx context.Context, status string, page, pageSize int) ([]*domain.Task, error)
    GetTaskByID(ctx context.Context, taskID string) (*domain.Task, error)
    UpdateInDb(ctx context.Context, taskID string, updates domain.TaskUpdate) error
    DeleteFromDb(ctx context.Context, taskID string) error
}
```

This allows for:
- ✅ Easy testing with mock implementations
- 🔄 Switching between different storage backends
- 🛠️ Implementation flexibility

## 🚀 Getting Started

### 📋 Prerequisites
- Go 1.25 or higher
- MongoDB instance (local or remote)

### 💿 Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd task-management
```

2. Install dependencies:
```bash
go mod tidy
```

3. Start MongoDB:
```bash
# Using Docker
docker run --name mongodb -p 27017:27017 -d mongo:latest

# Or start your local MongoDB instance
mongod
```

4. Run the service:
```bash
go run cmd/server/main.go
```

The service will start on port 8080.

### ⚙️ Configuration

The service currently connects to MongoDB at `mongodb://localhost:27017`. To change this, modify the connection string in `internal/repository/task_repository.go`.

## 📘 API Documentation

Base URL: `http://localhost:8080`

### 🔗 Endpoints

#### 📝 Create Task
```http
POST /tasks
Content-Type: application/json

{
    "title": "Complete project documentation",
    "description": "Write comprehensive README and API documentation"
}
```

**Response:**
```json
{
    "message": "Record Successfully Created",
    "task": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "Complete project documentation",
        "description": "Write comprehensive README and API documentation",
        "status": "pending",
        "created_at": "2023-12-07T10:30:00Z",
        "updated_at": "0001-01-01T00:00:00Z"
    }
}
```

#### 📋 List Tasks
```http
GET /tasks?status=pending&page=1&pageSize=10
```

**Query Parameters:**
- `status` (optional): Filter by task status (pending, in_progress, completed)
- `page` (optional): Page number (default: 1)
- `pageSize` (optional): Items per page (default: 10)

**Response:**
```json
{
    "data": [
        {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "title": "Complete project documentation",
            "description": "Write comprehensive README and API documentation",
            "status": "pending",
            "created_at": "2023-12-07T10:30:00Z",
            "updated_at": "0001-01-01T00:00:00Z"
        }
    ],
    "pagination": {
        "page": 1,
        "pageSize": 10
    }
}
```

#### 🔍 Get Task by ID
```http
GET /tasks/{id}
```

**Response:**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Complete project documentation",
    "description": "Write comprehensive README and API documentation",
    "status": "pending",
    "created_at": "2023-12-07T10:30:00Z",
    "updated_at": "0001-01-01T00:00:00Z"
}
```

#### 🔄 Update Task
```http
PUT /tasks/{id}
Content-Type: application/json

{
    "title": "Updated title",
    "description": "Updated description",
    "status": "in_progress"
}
```

**Response:** Returns the updated task object.

#### 🔄 Update Task Status Only
```http
PATCH /tasks/{id}
Content-Type: application/json

{
    "status": "completed"
}
```

**Response:** 200 OK (no body)

#### 🗑️ Delete Task
```http
DELETE /tasks/{id}
```

**Response:** 204 No Content

### ⚠️ Error Responses

#### 400 Bad Request
```json
{
    "error": "invalid status. Must be one of pending, in_progress, or completed"
}
```

#### 404 Not Found
```json
{
    "error": "task not found"
}
```

#### 500 Internal Server Error
```json
{
    "error": "internal server error description"
}
```

## 📈 Scalability

### 🔄 Horizontal Scaling

The service is designed for horizontal scaling:

1. **Stateless Design**: No server-side state is maintained
2. **Database Separation**: Data persistence is externalized to MongoDB
3. **Load Balancer Ready**: Multiple instances can run behind a load balancer
4. **Context-Aware**: All operations use context for timeout and cancellation

### 🛠️ Scaling Strategies

#### 🐳 Container Deployment
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o task-service cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/task-service .
CMD ["./task-service"]
```

#### ☸️ Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: task-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: task-service
  template:
    metadata:
      labels:
        app: task-service
    spec:
      containers:
      - name: task-service
        image: task-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: MONGO_URI
          value: "mongodb://mongodb-service:27017"
```

#### 🔄 Database Scaling
- **Read Replicas**: MongoDB read replicas for read-heavy workloads
- **Sharding**: Horizontal partitioning for very large datasets
- **Indexing**: Proper indexing on frequently queried fields (status, created_at)
