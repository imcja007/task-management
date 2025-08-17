package api

import (
	"encoding/json"
	"net/http"

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
	tasks, err := h.service.ListTasks(r.Context(), status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
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

//	Uncomment to insert random todo.

// func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {

// 	resp, err := http.Get("https://dummyjson.com/todos/random")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)

// 	var result map[string]interface{}
// 	if err := json.Unmarshal(body, &result); err != nil {
// 		log.Fatal(err)
// 	}
// 	task, err := h.service.CreateTask(r.Context(), result["todo"].(string), result["todo"].(string))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	response := map[string]interface{}{
// 		"message": "Record Successfully Created",
// 		"task":    task,
// 	}
// 	respondWithJSON(w, http.StatusCreated, response)
// }