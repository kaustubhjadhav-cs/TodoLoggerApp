package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// IST timezone (UTC+5:30)
var IST = time.FixedZone("IST", 5*60*60+30*60)

// GetTodayIST returns today's date in IST
func GetTodayIST() string {
	return time.Now().In(IST).Format("2006-01-02")
}

// GetYesterdayIST returns yesterday's date in IST
func GetYesterdayIST() string {
	return time.Now().In(IST).AddDate(0, 0, -1).Format("2006-01-02")
}

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./todo.db")
	if err != nil {
		return err
	}

	// Create tables
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		color TEXT NOT NULL DEFAULT '#58a6ff',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT DEFAULT '',
		created_date TEXT NOT NULL,
		assigned_date TEXT NOT NULL,
		completed_date TEXT,
		is_completed BOOLEAN DEFAULT FALSE,
		category_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
	);

	CREATE INDEX IF NOT EXISTS idx_assigned_date ON tasks(assigned_date);
	CREATE INDEX IF NOT EXISTS idx_completed_date ON tasks(completed_date);
	CREATE INDEX IF NOT EXISTS idx_category_id ON tasks(category_id);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	// Add category_id column if it doesn't exist (for existing databases)
	db.Exec(`ALTER TABLE tasks ADD COLUMN category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL`)

	// Insert default categories if none exist
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM categories`).Scan(&count)
	if err != nil || count == 0 {
		db.Exec(`INSERT OR IGNORE INTO categories (name, color) VALUES ('Work', '#58a6ff')`)
		db.Exec(`INSERT OR IGNORE INTO categories (name, color) VALUES ('Personal', '#3fb950')`)
		db.Exec(`INSERT OR IGNORE INTO categories (name, color) VALUES ('Misc', '#f0883e')`)
	}

	return nil
}

// Category CRUD operations

// CreateCategory creates a new category
func CreateCategory(name, color string) (*Category, error) {
	result, err := db.Exec(
		`INSERT INTO categories (name, color) VALUES (?, ?)`,
		name, color,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return GetCategoryByID(id)
}

// GetCategoryByID retrieves a category by ID
func GetCategoryByID(id int64) (*Category, error) {
	cat := &Category{}
	var createdAt string

	err := db.QueryRow(
		`SELECT id, name, color, created_at FROM categories WHERE id = ?`,
		id,
	).Scan(&cat.ID, &cat.Name, &cat.Color, &createdAt)

	if err != nil {
		return nil, err
	}

	cat.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return cat, nil
}

// GetAllCategories retrieves all categories with task counts
func GetAllCategories() ([]Category, error) {
	rows, err := db.Query(
		`SELECT c.id, c.name, c.color, c.created_at, 
		 (SELECT COUNT(*) FROM tasks WHERE category_id = c.id AND is_completed = FALSE) as task_count
		 FROM categories c ORDER BY c.name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		var createdAt string

		err := rows.Scan(&cat.ID, &cat.Name, &cat.Color, &createdAt, &cat.TaskCount)
		if err != nil {
			return nil, err
		}

		cat.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		categories = append(categories, cat)
	}

	return categories, nil
}

// UpdateCategory updates a category
func UpdateCategory(id int64, name, color string) (*Category, error) {
	_, err := db.Exec(
		`UPDATE categories SET name = ?, color = ? WHERE id = ?`,
		name, color, id,
	)
	if err != nil {
		return nil, err
	}

	return GetCategoryByID(id)
}

// DeleteCategory deletes a category
func DeleteCategory(id int64) error {
	_, err := db.Exec(`DELETE FROM categories WHERE id = ?`, id)
	return err
}

// CreateTask creates a new task
func CreateTask(title, description, date string) (*Task, error) {
	if date == "" {
		date = GetTodayIST()
	}

	result, err := db.Exec(
		`INSERT INTO tasks (title, description, created_date, assigned_date, is_completed) VALUES (?, ?, ?, ?, ?)`,
		title, description, date, date, false,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return GetTaskByID(id)
}

// GetTaskByID retrieves a task by ID
func GetTaskByID(id int64) (*Task, error) {
	task := &Task{}
	var completedDate sql.NullString
	var categoryID sql.NullInt64
	var createdAt, updatedAt string

	err := db.QueryRow(
		`SELECT id, title, description, created_date, assigned_date, completed_date, is_completed, category_id, created_at, updated_at FROM tasks WHERE id = ?`,
		id,
	).Scan(&task.ID, &task.Title, &task.Description, &task.CreatedDate, &task.AssignedDate, &completedDate, &task.IsCompleted, &categoryID, &createdAt, &updatedAt)

	if err != nil {
		return nil, err
	}

	if completedDate.Valid {
		task.CompletedDate = &completedDate.String
	}

	if categoryID.Valid {
		task.CategoryID = &categoryID.Int64
		task.Category, _ = GetCategoryByID(categoryID.Int64)
	}

	task.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	task.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	// Calculate drag days
	task.DragDays = CalculateBusinessDays(task.CreatedDate, task.AssignedDate)

	return task, nil
}

// GetTasksByDate retrieves all tasks for a specific date
func GetTasksByDate(date string) ([]Task, error) {
	rows, err := db.Query(
		`SELECT id, title, description, created_date, assigned_date, completed_date, is_completed, category_id, created_at, updated_at 
		 FROM tasks WHERE assigned_date = ? OR (completed_date = ? AND is_completed = TRUE)
		 ORDER BY is_completed ASC, created_at ASC`,
		date, date,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var completedDate sql.NullString
		var categoryID sql.NullInt64
		var createdAt, updatedAt string

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedDate, &task.AssignedDate, &completedDate, &task.IsCompleted, &categoryID, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		if completedDate.Valid {
			task.CompletedDate = &completedDate.String
		}

		if categoryID.Valid {
			task.CategoryID = &categoryID.Int64
			task.Category, _ = GetCategoryByID(categoryID.Int64)
		}

		task.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		task.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		task.DragDays = CalculateBusinessDays(task.CreatedDate, task.AssignedDate)

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetDailyLog retrieves the daily log for a specific date
func GetDailyLog(date string) (*DailyLog, error) {
	tasks, err := GetTasksByDate(date)
	if err != nil {
		return nil, err
	}

	log := &DailyLog{
		Date:  date,
		Tasks: tasks,
	}

	for _, task := range tasks {
		if task.IsCompleted {
			log.CompletedCount++
		} else {
			log.PendingCount++
		}
	}

	return log, nil
}

// GetAllDates retrieves all unique dates that have tasks
func GetAllDates() ([]string, error) {
	rows, err := db.Query(
		`SELECT DISTINCT date FROM (
			SELECT assigned_date as date FROM tasks
			UNION
			SELECT completed_date as date FROM tasks WHERE completed_date IS NOT NULL
		) ORDER BY date DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}

	return dates, nil
}

// HistorySummary represents what was accomplished on a specific date
type HistorySummary struct {
	Date           string `json:"date"`
	CompletedCount int    `json:"completed_count"`
	PendingCount   int    `json:"pending_count"` // Tasks assigned but not yet completed
}

// GetHistorySummaries retrieves completion stats for all dates
func GetHistorySummaries() ([]HistorySummary, error) {
	// Get all unique dates (both assigned and completed)
	dates, err := GetAllDates()
	if err != nil {
		return nil, err
	}

	var summaries []HistorySummary

	for _, date := range dates {
		summary := HistorySummary{Date: date}

		// Count tasks COMPLETED on this date
		err := db.QueryRow(
			`SELECT COUNT(*) FROM tasks WHERE completed_date = ? AND is_completed = TRUE`,
			date,
		).Scan(&summary.CompletedCount)
		if err != nil {
			return nil, err
		}

		// Count tasks ASSIGNED to this date that are still pending
		err = db.QueryRow(
			`SELECT COUNT(*) FROM tasks WHERE assigned_date = ? AND is_completed = FALSE`,
			date,
		).Scan(&summary.PendingCount)
		if err != nil {
			return nil, err
		}

		// Only include dates that have some activity
		if summary.CompletedCount > 0 || summary.PendingCount > 0 {
			summaries = append(summaries, summary)
		}
	}

	return summaries, nil
}

// UpdateTaskCompletion marks a task as completed or not completed
func UpdateTaskCompletion(id int64, isCompleted bool) (*Task, error) {
	var completedDate interface{}
	if isCompleted {
		completedDate = GetTodayIST()
	} else {
		completedDate = nil
	}

	_, err := db.Exec(
		`UPDATE tasks SET is_completed = ?, completed_date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		isCompleted, completedDate, id,
	)
	if err != nil {
		return nil, err
	}

	return GetTaskByID(id)
}

// RolloverTasks moves incomplete tasks from one date to another
func RolloverTasks(fromDate, toDate string) (int, error) {
	result, err := db.Exec(
		`UPDATE tasks SET assigned_date = ?, updated_at = CURRENT_TIMESTAMP WHERE assigned_date = ? AND is_completed = FALSE`,
		toDate, fromDate,
	)
	if err != nil {
		return 0, err
	}

	affected, _ := result.RowsAffected()
	return int(affected), nil
}

// RolloverAllPendingTasks moves ALL incomplete tasks from any past date to today
func RolloverAllPendingTasks(toDate string) (int, error) {
	result, err := db.Exec(
		`UPDATE tasks SET assigned_date = ?, updated_at = CURRENT_TIMESTAMP WHERE assigned_date < ? AND is_completed = FALSE`,
		toDate, toDate,
	)
	if err != nil {
		return 0, err
	}

	affected, _ := result.RowsAffected()
	return int(affected), nil
}

// DeleteTask deletes a task by ID
func DeleteTask(id int64) error {
	_, err := db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}

// UpdateTask updates a task's title and description
func UpdateTask(id int64, title, description string) (*Task, error) {
	_, err := db.Exec(
		`UPDATE tasks SET title = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		title, description, id,
	)
	if err != nil {
		return nil, err
	}

	return GetTaskByID(id)
}

// UpdateTaskCategory updates a task's category
func UpdateTaskCategory(id int64, categoryID *int64) (*Task, error) {
	_, err := db.Exec(
		`UPDATE tasks SET category_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		categoryID, id,
	)
	if err != nil {
		return nil, err
	}

	return GetTaskByID(id)
}

// GetTasksByCategory retrieves all incomplete tasks for a specific category
func GetTasksByCategory(categoryID int64) ([]Task, error) {
	rows, err := db.Query(
		`SELECT id, title, description, created_date, assigned_date, completed_date, is_completed, category_id, created_at, updated_at 
		 FROM tasks WHERE category_id = ? AND is_completed = FALSE
		 ORDER BY assigned_date ASC, created_at ASC`,
		categoryID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var completedDate sql.NullString
		var catID sql.NullInt64
		var createdAt, updatedAt string

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedDate, &task.AssignedDate, &completedDate, &task.IsCompleted, &catID, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		if completedDate.Valid {
			task.CompletedDate = &completedDate.String
		}

		if catID.Valid {
			task.CategoryID = &catID.Int64
			task.Category, _ = GetCategoryByID(catID.Int64)
		}

		task.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		task.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		task.DragDays = CalculateBusinessDays(task.CreatedDate, task.AssignedDate)

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// CalculateBusinessDays calculates the number of business days between two dates
func CalculateBusinessDays(startDate, endDate string) int {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return 0
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return 0
	}

	if start.Equal(end) || start.After(end) {
		return 0
	}

	businessDays := 0
	current := start

	for current.Before(end) {
		current = current.AddDate(0, 0, 1)
		// Skip weekends (Saturday = 6, Sunday = 0)
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			businessDays++
		}
	}

	return businessDays
}

// GetCompletedTasksForDate retrieves tasks that were completed on a specific date
func GetCompletedTasksForDate(date string) ([]Task, error) {
	rows, err := db.Query(
		`SELECT id, title, description, created_date, assigned_date, completed_date, is_completed, category_id, created_at, updated_at 
		 FROM tasks WHERE completed_date = ? AND is_completed = TRUE
		 ORDER BY created_at ASC`,
		date,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var completedDate sql.NullString
		var categoryID sql.NullInt64
		var createdAt, updatedAt string

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedDate, &task.AssignedDate, &completedDate, &task.IsCompleted, &categoryID, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		if completedDate.Valid {
			task.CompletedDate = &completedDate.String
		}

		if categoryID.Valid {
			task.CategoryID = &categoryID.Int64
			task.Category, _ = GetCategoryByID(categoryID.Int64)
		}

		task.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		task.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		task.DragDays = CalculateBusinessDays(task.CreatedDate, task.AssignedDate)

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetHistoricalLog retrieves the log of what was accomplished on a specific date
func GetHistoricalLog(date string) (*DailyLog, error) {
	// Get tasks that were completed on this date
	completedTasks, err := GetCompletedTasksForDate(date)
	if err != nil {
		return nil, err
	}

	// Get tasks that were assigned to this date but not completed (they would have been rolled over)
	rows, err := db.Query(
		`SELECT id, title, description, created_date, assigned_date, completed_date, is_completed, created_at, updated_at 
		 FROM tasks WHERE created_date <= ? AND (
			 (completed_date = ?) OR 
			 (assigned_date > ? AND created_date <= ?)
		 )
		 ORDER BY is_completed DESC, created_at ASC`,
		date, date, date, date,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying historical tasks: %v", err)
	}
	defer rows.Close()

	log := &DailyLog{
		Date:  date,
		Tasks: completedTasks,
	}

	log.CompletedCount = len(completedTasks)

	return log, nil
}
