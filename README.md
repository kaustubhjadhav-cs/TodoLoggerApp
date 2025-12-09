# Daily Task Tracker

A beautiful, feature-rich todo app with daily task tracking, automatic rollover of incomplete tasks, business day tracking, and optional category organization.

## Features

### Core Features
- **Daily Task Board**: View and manage tasks for any specific date
- **Task Completion Tracking**: Mark tasks as complete with visual feedback
- **Automatic Rollover**: Move all incomplete tasks from any past date to today with one click
- **Drag Day Tracking**: See how many business days (excluding weekends) a task has been pending
- **Historical Logs**: Browse and view what was accomplished on each day
- **Progress Statistics**: Real-time stats showing completed, pending, total, and dragged tasks

### Theme & UI
- **Light/Dark Mode**: Toggle between themes with the 🌙/☀️ button
- **Theme Persistence**: Your preference is saved and remembered
- **Modern UI**: Beautiful design with smooth animations and responsive layout
- **Custom Modals**: Clean confirmation dialogs (no browser alerts)

### Categories (Optional)
- **Organize Tasks**: Assign categories like Work, Personal, Misc to tasks
- **Color Coded**: Each category has a customizable color
- **Filter by Category**: View tasks by category in the sidebar or dropdown
- **Manage Categories**: Create, rename, and delete categories in Settings
- **Toggle Feature**: Enable/disable categories entirely in Settings

### Settings & Customization
- **Show Creation Date**: Display when tasks were originally created
- **Show Drag Days**: Display how many days tasks have been pending
- **Show Assigned Date**: Display the date tasks are assigned to
- **Enable Categories**: Toggle the categories feature on/off

### Timezone
- **IST Support**: All dates and times use Indian Standard Time (Asia/Kolkata)

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

### Building a Binary

For faster startup without compilation:
```bash
go build -o todoapp .
./todoapp
```

## Usage

### Adding Tasks
- Type your task in the input field and press **Enter** or click "Add Task"
- Tasks are automatically assigned to the currently selected date

### Completing Tasks
- Click the checkbox next to a task to mark it as complete
- Completed tasks show with a strikethrough and green checkmark

### Navigating Dates
- Use the **← Prev** / **Next →** buttons to navigate between dates
- Use the date picker to jump to any specific date
- Click **Today** to quickly return to today's date

### Rolling Over Tasks
- Click **"Rollover Pending"** to move ALL incomplete tasks from any past date to today
- Tasks retain their creation date for accurate drag day tracking

### Viewing History
- The **Historical Logs** sidebar shows dates with task activity
- Counts show tasks **completed ON that day** (not just assigned)
- Click any date to view details in a modal

### Using Categories
1. Go to **Settings** (⚙️ button)
2. Enable **"Enable Categories"**
3. Hover over any task → click 🏷️ to assign a category
4. Use the filter dropdown or sidebar to view tasks by category

### Managing Categories
1. Go to **Settings** (⚙️ button)
2. Scroll to **"Manage Categories"** section
3. Click color picker to change category color
4. Edit name and press **Enter** to rename (confirmation modal appears)
5. Click 🗑️ to delete a category

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `/` | Focus the task input field |
| `t` | Toggle light/dark theme |
| `Enter` | Confirm action in modals |
| `Escape` | Close any open modal |

## API Endpoints

### Tasks

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/tasks?date=YYYY-MM-DD` | Get tasks for a specific date |
| POST | `/api/tasks` | Create a new task |
| PUT | `/api/tasks/{id}` | Update a task |
| DELETE | `/api/tasks/{id}` | Delete a task |
| PUT | `/api/tasks/{id}/complete` | Toggle task completion |
| PUT | `/api/tasks/{id}/category` | Update task's category |

### Daily Logs & History

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/daily-log?date=YYYY-MM-DD` | Get daily log with stats |
| GET | `/api/dates` | Get all dates with tasks |
| GET | `/api/history-summaries` | Get completion stats for all dates |

### Rollover

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/rollover` | Rollover tasks between specific dates |
| POST | `/api/rollover-all` | Rollover all past incomplete tasks to today |
| POST | `/api/auto-rollover` | Auto rollover from yesterday to today |

### Categories

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/categories` | Get all categories with task counts |
| POST | `/api/categories` | Create a new category |
| PUT | `/api/categories/{id}` | Update a category |
| DELETE | `/api/categories/{id}` | Delete a category |
| GET | `/api/categories/{id}/tasks` | Get tasks for a category |

## Drag Day Calculation

The app calculates "drag days" - the number of **business days** a task has been pending since its creation:
- Only counts weekdays (Monday-Friday)
- Excludes Saturday and Sunday
- Shows **orange** warning when dragging begins (1-2 days)
- Shows **red** critical warning when dragged 3+ days

## Project Structure

```
TodoApp/
├── main.go           # Application entry point and HTTP server
├── models.go         # Data structures (Task, Category, DailyLog)
├── database.go       # SQLite database operations & IST timezone
├── handlers.go       # HTTP request handlers
├── go.mod            # Go module dependencies
├── go.sum            # Dependency checksums
├── todo.db           # SQLite database (created on first run)
├── static/
│   └── index.html    # Frontend application (HTML, CSS, JS)
└── README.md         # This file
```

## Data Storage

- **Database**: SQLite (file-based, no external server needed)
- **Location**: `todo.db` in the project directory
- **Persistence**: All data persists across server restarts
- **Settings**: User preferences stored in browser's localStorage

## Browser Support

Modern browsers with ES6+ support:
- Chrome (recommended)
- Firefox
- Safari
- Edge

## License

MIT License - feel free to use and modify as needed.
