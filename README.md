# ✅ Task Management System

A RESTful microservice for managing tasks built with Go, MongoDB, and following microservices architecture principles.

## 📑 Table of Contents
- [🔍 Problem Breakdown](#-problem-breakdown)
- [🧩 Design Decisions](#-design-decisions)
- [🏗️ Architecture Overview](#️-architecture-overview)
- [🧠 Microservices Concepts](#-microservices-concepts)
- [🚀 Getting Started](#-getting-started)
- [🐳 Docker Setup](#-docker-setup)
- [📘 API Documentation](#-api-documentation)
- [📈 Scalability](#-scalability)
- [🔄 Inter-Service Communication](#-inter-service-communication)


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
- **Docker** 🐳: Containerization for easy deployment and scaling

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
├── Dockerfile           # Container configuration
├── docker-compose.yml   # Multi-service orchestration
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
- **Option 1 (Docker)**: Docker and Docker Compose
- **Option 2 (Manual)**: Go 1.25+ and MongoDB instance

### 💿 Installation

#### Option 1: Docker Setup (Recommended)
See the [🐳 Docker Setup](#-docker-setup) section below for the easiest way to get started.

#### Option 2: Manual Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd task-management
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set environment variables:
```bash
export MONGO_URI="mongodb://localhost:27017/taskdb"
```

4. Start MongoDB:
```bash
# Using Docker
docker run --name mongodb -p 27017:27017 -d mongo:latest

# Or start your local MongoDB instance
mongod
```

5. Run the service:
```bash
go run cmd/server/main.go
```

The service will start on port 8080.

## 🐳 Docker Setup

The easiest way to run the Task Management System is using Docker Compose, which will automatically set up both the application and MongoDB.

### 🚀 Quick Start with Docker

1. **Clone the repository:**
```bash
git clone https://github.com/imcja007/task-management.git
cd task-management
```

2. **Start the services:**
```bash
docker-compose up -d
```

This command will:
- Build the Go application container
- Start a MongoDB container with authentication
- Set up a network for communication between services
- Expose the API on port 8080

3. **Verify the services are running:**
```bash
docker-compose ps
```

You should see output similar to:
```
NAME                   COMMAND                  SERVICE      STATUS        PORTS
task-service-app       "go run cmd/server/m…"   task-service running       0.0.0.0:8080->8080/tcp
task-service-mongodb   "docker-entrypoint.s…"   mongodb      running       0.0.0.0:27017->27017/tcp
```

4. **Test the API:**
```bash
curl http://localhost:8080/tasks
```

5. **View logs (optional):**
```bash
# View all logs
docker-compose logs

# View only app logs
docker-compose logs task-service

# Follow logs in real-time
docker-compose logs -f
```

### 🛑 Stop the Services

```bash
# Stop services but keep data
docker-compose down

# Stop services and remove data volumes
docker-compose down -v
```

### 🔧 Docker Configuration Details

The `docker-compose.yml` file includes:

**MongoDB Service:**
- Image: `mongo:7.0`
- Port: `27017`
- Authentication: username `admin`, password `password123`
- Database: `taskdb`
- Persistent data storage with volumes

**Task Service:**
- Built from local Dockerfile
- Port: `8080`
- Environment variables automatically configured
- Waits for MongoDB to be ready
- Connected via internal Docker network

### 🐳 Manual Docker Commands

If you prefer to build and run containers manually:

```bash
# Build the application image
docker build -t task-management .

# Create a network
docker network create task-network

# Run MongoDB
docker run -d --name mongodb \
  --network task-network \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=password123 \
  -e MONGO_INITDB_DATABASE=taskdb \
  -p 27017:27017 \
  mongo:7.0

# Run the application
docker run -d --name task-service \
  --network task-network \
  -e MONGO_URI="mongodb://admin:password123@mongodb:27017/taskdb?authSource=admin" \
  -p 8080:8080 \
  task-management
```

### ⚙️ Configuration

The service uses environment variables for configuration:
- `MONGO_URI`: MongoDB connection string (required)
- `PORT`: Server port (default: 8080)
- `ENV`: Environment (development/production)

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

#### 🎲 Create Random Task
```http
POST /random-tasks
```

**Description:** Creates a random task by fetching data from an external API (dummyjson.com). This endpoint is useful for testing and populating the database with sample data.

**Response:**
```json
{
    "message": "Record Successfully Created",
    "task": {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "title": "Random task from external API",
        "description": "Random task from external API",
        "status": "pending",
        "created_at": "2023-12-07T10:35:00Z",
        "updated_at": "0001-01-01T00:00:00Z"
    }
}
```

#### 📋 List Tasks
```http
GET /tasks?status=pending&page=1&pageSize=10
```

**Query Parameters:**
- `status` (optional): Filter by task status (`pending`, `in_progress`, `completed`)
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

**Response:**
```json
{
    "message": "Record Successfully updated"
}
```

#### 🗑️ Delete Task
```http
DELETE /tasks/{id}
```

**Response:**
```json
{
    "message": "Record Successfully deleted"
}
```

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

### 🧪 API Testing Examples

You can test the API using curl commands:

```bash
# Create a new task
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Task", "description": "This is a test task"}'

# Create a random task
curl -X POST http://localhost:8080/random-tasks

# List all tasks
curl http://localhost:8080/tasks

# List tasks with filtering and pagination
curl "http://localhost:8080/tasks?status=pending&page=1&pageSize=5"

# Get a specific task (replace {id} with actual task ID)
curl http://localhost:8080/tasks/{id}

# Update task status
curl -X PATCH http://localhost:8080/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'

# Delete a task
curl -X DELETE http://localhost:8080/tasks/{id}
```

## 📈 Scalability

### 🔄 Horizontal Scaling

The service is designed for horizontal scaling:

1. **Stateless Design**: No server-side state is maintained
2. **Database Separation**: Data persistence is externalized to MongoDB
3. **Load Balancer Ready**: Multiple instances can run behind a load balancer
4. **Context-Aware**: All operations use context for timeout and cancellation
5. **Containerized**: Docker containers make scaling easy

### 🛠️ Scaling Strategies

#### 🐳 Container Orchestration with Kubernetes

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
        image: task-management:latest
        ports:
        - containerPort: 8080
        env:
        - name: MONGO_URI
          value: "mongodb://admin:password123@mongodb-service:27017/taskdb?authSource=admin"
---
apiVersion: v1
kind: Service
metadata:
  name: task-service
spec:
  selector:
    app: task-service
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

#### 🔄 Database Scaling
- **Read Replicas**: MongoDB read replicas for read-heavy workloads
- **Sharding**: Horizontal partitioning for very large datasets
- **Indexing**: Proper indexing on frequently queried fields (status, created_at)

#### 📊 Load Balancing
```yaml
# docker-compose.yml with load balancer
version: '3.8'
services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - task-service-1
      - task-service-2
      - task-service-3

  task-service-1:
    build: .
    environment:
      - MONGO_URI=mongodb://admin:password123@mongodb:27017/taskdb?authSource=admin

  task-service-2:
    build: .
    environment:
      - MONGO_URI=mongodb://admin:password123@mongodb:27017/taskdb?authSource=admin

  task-service-3:
    build: .
    environment:
      - MONGO_URI=mongodb://admin:password123@mongodb:27017/taskdb?authSource=admin
```

## 🔄 Inter-Service Communication

### 👥 Adding a User Service

When extending the system with additional microservices (e.g., User Service), several communication patterns can be employed:

#### 1️⃣ REST API Communication
```go
// User service client
type UserServiceClient struct {
    baseURL string
    client  *http.Client
}

func (c *UserServiceClient) GetUser(userID string) (*User, error) {
    resp, err := c.client.Get(fmt.Sprintf("%s/users/%s", c.baseURL, userID))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }
    return &user, nil
}
```

#### 2️⃣ gRPC Communication
```protobuf
syntax = "proto3";

service UserService {
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
}

message GetUserRequest {
    string user_id = 1;
}

message GetUserResponse {
    User user = 1;
}
```

#### 3️⃣ Event-Driven Communication
```go
// Task created event
type TaskCreatedEvent struct {
    TaskID    string    `json:"task_id"`
    UserID    string    `json:"user_id"`
    Title     string    `json:"title"`
    CreatedAt time.Time `json:"created_at"`
}

// Publish event after task creation
func (s *TaskService) CreateTask(ctx context.Context, userID, title, description string) (*domain.Task, error) {
    task := &domain.Task{
        ID:          uuid.New().String(),
        UserID:      userID,
        Title:       title,
        Description: description,
        Status:      "pending",
        CreatedAt:   time.Now(),
    }
    
    if err := s.repo.Create(ctx, task); err != nil {
        return nil, err
    }
    
    // Publish event
    event := TaskCreatedEvent{
        TaskID:    task.ID,
        UserID:    task.UserID,
        Title:     task.Title,
        CreatedAt: task.CreatedAt,
    }
    s.eventPublisher.Publish("task.created", event)
    
    return task, nil
}
```

### 🔄 Communication Patterns

1. **Synchronous (REST/gRPC)** ⚡: For real-time data requirements
2. **Asynchronous (Message Queues)** 📨: For eventual consistency and loose coupling
3. **Event Sourcing** 📊: For audit trails and complex business workflows

### 🔍 Service Discovery

For production deployments, implement service discovery:

```go
// Service registry interface
type ServiceRegistry interface {
    Register(serviceName, address string) error
    Discover(serviceName string) ([]string, error)
    Health(serviceName string) error
}

// Consul implementation
type ConsulRegistry struct {
    client *consulapi.Client
}
```

## 🎯 Quick Reference

### Essential Commands
```bash
# Start with Docker (recommended)
docker-compose up -d

# Manual start
export MONGO_URI="mongodb://localhost:27017/taskdb"
go run cmd/server/main.go

# Test the API
curl http://localhost:8080/tasks

# Stop Docker services
docker-compose down
```

### Key Features
- ✅ Full CRUD operations
- 🔄 Status management (pending, in_progress, completed)
- 📄 Pagination and filtering
- 🎲 Random task generation
- 🐳 Docker containerization
- 🏗️ Microservices architecture
- 📊 Production-ready scaling options
