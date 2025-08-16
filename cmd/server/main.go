package main

import (
	"log"
	"net/http"

	"task-management/api"
	"task-management/internal/repository"
	"task-management/internal/service"

	"github.com/gorilla/mux"
)

func main() {

	repo := repository.NewInMemoryTaskRepository()

	taskService := service.NewTaskService(repo)


	taskHandler := api.NewTaskHandler(taskService)

	router := mux.NewRouter()
	taskHandler.RegisterRoutes(router)

	// Start server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
