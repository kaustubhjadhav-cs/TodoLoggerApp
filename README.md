# Daily Task Tracker

A beautiful, minimalist todo app with daily task tracking, automatic rollover of incomplete tasks, and business day tracking for dragged tasks.

## Features

- **Daily Task Board**: View and manage tasks for any specific date
- **Task Completion Tracking**: Mark tasks as complete with visual feedback
- **Automatic Rollover**: Incomplete tasks from the previous day automatically move to today
- **Drag Day Tracking**: See how many business days (excluding weekends) a task has been pending
- **Historical Logs**: Browse and view completed tasks from any past date
- **Progress Statistics**: Real-time stats showing completed, pending, total, and dragged tasks
- **Modern UI**: Dark theme with smooth animations and responsive design

## Getting Started

### Prerequisites

- Go 1.21 or higher
- SQLite3 (CGO enabled)

### Installation

1. Clone or navigate to the project directory:
   ```bash
   cd TodoApp
   ```

2. Download dependencies:
   ```bash
   go mod tidy
   ```

3. Run the application:
   ```bash
   go run .
   ```

4. Open your browser and navigate to:
   ```
   http://localhost:8080
   ```

## Usage

### Adding Tasks
- Type your task in the input field and press Enter or click "Add Task"
- Tasks are automatically assigned to the currently selected date

### Completing Tasks
- Click the checkbox next to a task to mark it as complete
- Completed tasks show with a strikethrough and green checkmark

### Navigating Dates
- Use the date picker or arrow buttons to navigate between dates
- Click "Today" to quickly jump to today's date

### Rolling Over Tasks
- Click "Rollover Pending" to move incomplete tasks from yesterday to today
- This happens automatically for the previous day's tasks

### Viewing History
- The sidebar shows a list of dates with task activity
- Click on any date to view what was accomplished that day

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/tasks?date=YYYY-MM-DD` | Get tasks for a specific date |
| POST | `/api/tasks` | Create a new task |
| PUT | `/api/tasks/{id}` | Update a task |
| DELETE | `/api/tasks/{id}` | Delete a task |
| PUT | `/api/tasks/{id}/complete` | Toggle task completion |
| GET | `/api/daily-log?date=YYYY-MM-DD` | Get daily log with stats |
| GET | `/api/dates` | Get all dates with tasks |
| POST | `/api/rollover` | Rollover incomplete tasks between dates |
| POST | `/api/auto-rollover` | Auto rollover from yesterday to today |

## Drag Day Calculation

The app calculates "drag days" - the number of business days a task has been pending since its creation:
- Only counts weekdays (Monday-Friday)
- Excludes Saturday and Sunday
- Shows visual warning (orange) when dragging begins
- Shows critical warning (red) when dragged 3+ days

## Project Structure

```
TodoApp/
├── main.go           # Application entry point and HTTP server
├── models.go         # Data structures
├── database.go       # SQLite database operations
├── handlers.go       # HTTP request handlers
├── go.mod            # Go module dependencies
├── todo.db           # SQLite database (created on first run)
├── static/
│   └── index.html    # Frontend application
└── README.md         # This file
```

## Keyboard Shortcuts

- `/` - Focus the task input field
- `Escape` - Close any open modal
- `Enter` - Add task (when input is focused)

## License

MIT License - feel free to use and modify as needed.

