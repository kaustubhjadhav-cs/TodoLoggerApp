package main

import "time"

// Task represents a todo item
type Task struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	CreatedDate   string    `json:"created_date"`   // Date when task was first created
	AssignedDate  string    `json:"assigned_date"`  // Current date the task is assigned to
	CompletedDate *string   `json:"completed_date"` // Date when task was completed (nil if not completed)
	IsCompleted   bool      `json:"is_completed"`
	DragDays      int       `json:"drag_days"` // Business days the task has been dragged
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// DailyLog represents tasks for a specific day
type DailyLog struct {
	Date           string `json:"date"`
	Tasks          []Task `json:"tasks"`
	CompletedCount int    `json:"completed_count"`
	PendingCount   int    `json:"pending_count"`
}

// TaskRequest is used for creating/updating tasks
type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"` // Optional: defaults to today
}

// CompleteTaskRequest marks a task as complete
type CompleteTaskRequest struct {
	IsCompleted bool `json:"is_completed"`
}

// RolloverRequest handles end of day rollover
type RolloverRequest struct {
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
}
