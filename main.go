package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Initialize database
	if err := initDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Create a new mux
	mux := http.NewServeMux()

	// Serve static files
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			http.ServeFile(w, r, "static/index.html")
			return
		}
		http.ServeFile(w, r, "static"+r.URL.Path)
	})

	// API Routes
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGetTasks(w, r)
		case "POST":
			HandleCreateTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")

		if strings.HasSuffix(path, "/complete") {
			if r.Method == "PUT" {
				HandleUpdateTaskCompletion(w, r)
				return
			}
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		switch r.Method {
		case "GET":
			HandleGetTask(w, r)
		case "PUT":
			HandleUpdateTask(w, r)
		case "DELETE":
			HandleDeleteTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/daily-log", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			HandleGetDailyLog(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/dates", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			HandleGetAllDates(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/history-summaries", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			HandleGetHistorySummaries(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/rollover", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			HandleRollover(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/auto-rollover", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			HandleAutoRollover(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/rollover-all", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			HandleRolloverAll(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/historical-log", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			HandleGetHistoricalLog(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Apply CORS middleware
	handler := corsMiddleware(mux)

	fmt.Println("🚀 Todo App server starting on http://localhost:8080")
	fmt.Println("📝 Open your browser to http://localhost:8080 to use the app")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
