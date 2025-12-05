package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Response helpers
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// HandleCreateTask creates a new task
func HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		respondError(w, http.StatusBadRequest, "Title is required")
		return
	}

	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}

	task, err := CreateTask(req.Title, req.Description, req.Date)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, task)
}

// HandleGetTasks gets all tasks for a specific date
func HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	tasks, err := GetTasksByDate(date)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if tasks == nil {
		tasks = []Task{}
	}

	respondJSON(w, http.StatusOK, tasks)
}

// HandleGetDailyLog gets the daily log for a specific date
func HandleGetDailyLog(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	log, err := GetDailyLog(date)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if log.Tasks == nil {
		log.Tasks = []Task{}
	}

	respondJSON(w, http.StatusOK, log)
}

// HandleGetAllDates gets all dates that have tasks
func HandleGetAllDates(w http.ResponseWriter, r *http.Request) {
	dates, err := GetAllDates()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if dates == nil {
		dates = []string{}
	}

	respondJSON(w, http.StatusOK, dates)
}

// HandleUpdateTaskCompletion updates the completion status of a task
func HandleUpdateTaskCompletion(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path: /api/tasks/{id}/complete
	path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	path = strings.TrimSuffix(path, "/complete")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req CompleteTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task, err := UpdateTaskCompletion(id, req.IsCompleted)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, task)
}

// HandleRollover rolls over incomplete tasks to the next date
func HandleRollover(w http.ResponseWriter, r *http.Request) {
	var req RolloverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.FromDate == "" || req.ToDate == "" {
		respondError(w, http.StatusBadRequest, "Both from_date and to_date are required")
		return
	}

	count, err := RolloverTasks(req.FromDate, req.ToDate)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":     "Tasks rolled over successfully",
		"tasks_moved": count,
		"from_date":   req.FromDate,
		"to_date":     req.ToDate,
	})
}

// HandleDeleteTask deletes a task
func HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path: /api/tasks/{id}
	path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	if err := DeleteTask(id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

// HandleUpdateTask updates a task's title and description
func HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path: /api/tasks/{id}
	path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task, err := UpdateTask(id, req.Title, req.Description)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, task)
}

// HandleGetHistoricalLog gets the historical log for a specific date
func HandleGetHistoricalLog(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		respondError(w, http.StatusBadRequest, "Date is required")
		return
	}

	log, err := GetHistoricalLog(date)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if log.Tasks == nil {
		log.Tasks = []Task{}
	}

	respondJSON(w, http.StatusOK, log)
}

// HandleGetTask gets a single task by ID
func HandleGetTask(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path: /api/tasks/{id}
	path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := GetTaskByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

// HandleAutoRollover automatically rolls over incomplete tasks from yesterday to today
func HandleAutoRollover(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	count, err := RolloverTasks(yesterday, today)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":     "Auto rollover completed",
		"tasks_moved": count,
		"from_date":   yesterday,
		"to_date":     today,
	})
}

// HandleRolloverAll rolls over ALL incomplete tasks from any past date to today
func HandleRolloverAll(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format("2006-01-02")

	count, err := RolloverAllPendingTasks(today)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":     "All pending tasks rolled over to today",
		"tasks_moved": count,
		"to_date":     today,
	})
}
