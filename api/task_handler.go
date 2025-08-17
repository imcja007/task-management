package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"task-management/internal/domain"
	"task-management/internal/service"

	"github.com/gorilla/mux"
)

// TaskHandler handles HTTP requests for tasks
type TaskHandler struct {
	service *service.TaskService
}

// creates a new instance
func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/tasks", h.CreateTask).Methods("POST")
	router.HandleFunc("/random-tasks", h.CreateRandomTask).Methods("POST")
	router.HandleFunc("/tasks", h.ListTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", h.ListTaskByID).Methods("GET")
	router.HandleFunc("/tasks/{id}", h.UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", h.DeleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", h.UpdateTaskStatus).Methods("PATCH")
}


func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.service.CreateTask(r.Context(), req.Title, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Record Successfully Created",
		"task":    task,
	}
	respondWithJSON(w, http.StatusCreated, response)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	
	// Parse pagination parameters
	page := 1 // default page
	pageSize := 10 // default page size
	
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	tasks, err := h.service.ListTasks(r.Context(), status, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data": tasks,
		"pagination": map[string]interface{}{
			"page": page,
			"pageSize": pageSize,
		},
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaskHandler) ListTaskByID(w http.ResponseWriter, r *http.Request) {
	taskId := mux.Vars(r)["id"]
	task, err := h.service.GetTaskByIDFromDB(r.Context(), taskId)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	respondWithJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updates domain.TaskUpdate
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateTask(r.Context(), id, updates); err != nil {
		if err == domain.ErrTaskNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the updated task to return in response
	task, err := h.service.GetTaskByIDFromDB(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

// DeleteTask handles task deletion requests
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		if err == domain.ErrTaskNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (h *TaskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateTaskStatus(r.Context(), id, req.Status); err != nil {
		if err == domain.ErrTaskNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err == domain.ErrInvalidNewStatus {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

//	Uncomment to insert random tasks.

func (h *TaskHandler) CreateRandomTask(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get("https://dummyjson.com/todos/random")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}
	task, err := h.service.CreateTask(r.Context(), result["todo"].(string), result["todo"].(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Record Successfully Created",
		"task":    task,
	}
	respondWithJSON(w, http.StatusCreated, response)
}


//	create 20 tasks
// func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
// 	var tasks []*domain.Task
// 	for i := 0; i < 20; i++ {
// 		resp, err := http.Get("https://dummyjson.com/todos/random")
// 		if err != nil {
// 			log.Printf("Error fetching todo #%d: %v", i+1, err)
// 			http.Error(w, "Error fetching todos", http.StatusInternalServerError)
// 			return
// 		}	
// 		body, err := io.ReadAll(resp.Body)
// 		resp.Body.Close()	
// 		if err != nil {
// 			log.Printf("Error reading response body for todo #%d: %v", i+1, err)
// 			http.Error(w, "Error reading response", http.StatusInternalServerError)
// 			return
// 		}
// 		var result map[string]interface{}
// 		if err := json.Unmarshal(body, &result); err != nil {
// 			log.Printf("Error unmarshaling JSON for todo #%d: %v", i+1, err)
// 			http.Error(w, "Error parsing response", http.StatusInternalServerError)
// 			return
// 		}
// 		task, err := h.service.CreateTask(r.Context(), result["todo"].(string), result["todo"].(string))
// 		if err != nil {
// 			log.Printf("Error creating todo #%d: %v", i+1, err)
// 			http.Error(w, "Error creating tasks", http.StatusInternalServerError)
// 			return
// 		}

// 		tasks = append(tasks, task)
// 	}

// 	response := map[string]interface{}{
// 		"message": "Successfully Created 20 Records",
// 		"tasks":   tasks,
// 	}
// 	respondWithJSON(w, http.StatusCreated, response)
// }