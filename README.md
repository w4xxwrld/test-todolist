# Todo List Desktop Application

A cross-platform desktop Todo List application built with Wails v2 (Go + JavaScript). Features task management with persistent storage and responsive UI with theme support.

## Features Implemented

**Basic Features:**
- [x] User interface with input field and buttons
- [x] Add tasks with validation
- [x] Delete tasks
- [x] Mark tasks complete/incomplete
- [x] Save/load state with file persistence
- [x] Filter and sort tasks

**Bonus Features:**
- [x] Responsive layout
- [x] Dark/light theme toggle
- [x] Due dates and priority levels
- [x] Delete confirmation modal
- [x] Separate completed tasks section
- [x] PostgreSQL database integration
- [x] Advanced filtering by date ranges
- [x] Priority-based sorting

## How to Launch

### Prerequisites
- Go 1.19+
- Node.js 16+

### Setup
```bash
# Install Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Clone and setup project
git clone <repository-url>
cd test-todolist
go mod tidy
cd frontend && npm install && cd ..
```

### Run Development Server
```bash
$(go env GOPATH)/bin/wails dev
```

### Build Production Version
```bash
$(go env GOPATH)/bin/wails build
```

The build of the app will be in `build/bin/` directory.

### Alternative Launch Commands
If the above doesn't work, try:
```bash
# Option 1: Add to PATH temporarily
export PATH=$PATH:$(go env GOPATH)/bin && wails dev

# Option 2: Use full path directly
/Users/$(whoami)/go/bin/wails dev
```

## Data Storage

By default, tasks are saved to `~/.todolist/tasks.json`

For PostgreSQL support, set environment variable:
```bash
export POSTGRES_CONNECTION_STRING="postgres://user:pass@localhost/dbname?sslmode=disable"
```

## Project Structure
```
├── internal/           # Backend (Go)
│   ├── domain/         # Business entities
│   ├── repository/     # Data storage (memory/file/postgres)
│   ├── service/        # Business logic
│   └── usecase/        # Application layer
├── frontend/           # Frontend (HTML/CSS/JS)
├── migrations/         # Database schema
└── build/             # Build output
```

## Architecture
- Clean Architecture with Domain/Repository/Service/UseCase layers
- Multiple storage backends with automatic fallback
- Modern responsive frontend with vanilla JavaScript
- Cross-platform desktop application using Wails framework
