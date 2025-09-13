# 📄 [Download my CV](screenshots/Screenshot%202025-09-13%20at%209.10.57%E2%80%AFPM.png)

# Todo List Desktop Application

A cross-platform desktop Todo List application built with Wails v2 (Go + JavaScript). Features task management with persistent storage and responsive UI with theme support.

## Link for the demo
- https://www.loom.com/share/9b2c95188b844a1a8b4077567f19dac5

## Linkedin
https://www.linkedin.com/in/shakhnazar-mussabekov/

## Screenshots

Below are some screenshots of the Todo List app:

![Main UI](screenshots/Screenshot%202025-09-13%20at%209.10.57%E2%80%AFPM.png)
![Add Task](screenshots/Screenshot%202025-09-13%20at%209.40.01%E2%80%AFPM.png)
![Delete Confirmation](screenshots/Screenshot%202025-09-13%20at%209.40.18%E2%80%AFPM.png)
![Completed Tasks](screenshots/Screenshot%202025-09-13%20at%209.40.38%E2%80%AFPM.png)
![Dark Theme](screenshots/Screenshot%202025-09-13%20at%209.40.59%E2%80%AFPM.png)

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
