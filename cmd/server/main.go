package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if not specified
	}

	// Start server
	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatal(err)
	}
}
