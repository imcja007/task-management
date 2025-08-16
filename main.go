package main

import (
    "encoding/json"
    "log"
    "net/http"
)

type Response struct {
    Message string `json:"message"`
}

func main() {
    // handler for GET requests
    http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
        // Only allow GET requests
        if r.Method != http.MethodGet {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        response := Response{
            Message: "Hello from the API!",
        }

        w.Header().Set("Content-Type", "application/json")
        
        json.NewEncoder(w).Encode(response)
    })

    // Start the server
    log.Println("Server starting on http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
