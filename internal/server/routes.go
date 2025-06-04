package server

import (
	"context"
	"fx-todo-api/internal/tasks"
	"log"
	"net/http"

	"go.uber.org/fx"
)

// CORSMiddleware adds CORS headers and handles preflight requests.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// NewMux creates the HTTP request multiplexer and wraps it with middleware.
func NewMux(handler *tasks.TaskHandler) http.Handler {
	mux := http.NewServeMux()

	// Handler for /task (LIST and CREATE)
	// mux.HandleFunc("GET /task", handler.ListTask)
	mux.HandleFunc("POST /task", handler.CreateTask)

	// Handler for /task/ (UPDATE and DELETE specific tasks)
	mux.HandleFunc("UPDATE /task", handler.UpdateTask)
	mux.HandleFunc("DELETE /task", handler.DeleteTask)

	// Default handler for unknown routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	// Wrap the mux with CORS middleware
	return CORSMiddleware(mux)
}

// StartServer starts the HTTP server with lifecycle management.
func StartServer(lc fx.Lifecycle, handler *tasks.TaskHandler) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Println("Starting server on :8080")
				if err := http.ListenAndServe(":8080", NewMux(handler)); err != nil {
					log.Fatalf("Server failed: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping server")
			return nil
		},
	})
}
