package main

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"todo-list/internal/domain"
	"todo-list/internal/repository"
	"todo-list/internal/service"
	"todo-list/internal/usecase"
)

type App struct {
	ctx         context.Context
	taskUseCase *usecase.TaskUseCase
}

func NewApp() *App {
	var taskRepo repository.TaskRepository

	pgConnStr := os.Getenv("POSTGRES_CONNECTION_STRING")
	if pgConnStr != "" {
		pgRepo, err := repository.NewPostgresTaskRepository(pgConnStr)
		if err == nil {
			taskRepo = pgRepo
		} else {
			println("Failed to connect to PostgreSQL:", err.Error())
		}
	}

	if taskRepo == nil {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "."
		}

		dataDir := filepath.Join(homeDir, ".todolist")
		fileRepo, err := repository.NewFileTaskRepository(dataDir)
		if err != nil {
			taskRepo = repository.NewMemoryTaskRepository()
		} else {
			taskRepo = fileRepo
		}
	}

	taskService := service.NewTaskService(taskRepo)
	taskUseCase := usecase.NewTaskUseCase(taskService)

	return &App{
		taskUseCase: taskUseCase,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) CreateTask(title, description string) (*domain.Task, error) {
	req := usecase.CreateTaskRequest{
		Title:       title,
		Description: description,
	}
	return a.taskUseCase.CreateTask(a.ctx, req)
}

func (a *App) CreateTaskWithDetails(title, description, priority string, dueDate *time.Time) (*domain.Task, error) {
	req := usecase.CreateTaskRequest{
		Title:       title,
		Description: description,
		Priority:    priority,
		DueDate:     dueDate,
	}
	return a.taskUseCase.CreateTask(a.ctx, req)
}

func (a *App) GetAllTasks() ([]*domain.Task, error) {
	filter := usecase.TaskFilter{Status: "all"}
	sort := usecase.TaskSort{Field: "created", Order: "desc"}
	return a.taskUseCase.GetFilteredAndSortedTasks(a.ctx, filter, sort)
}

func (a *App) GetFilteredTasks(status, priority, dateType, sortField, sortOrder string) ([]*domain.Task, error) {
	filter := usecase.TaskFilter{
		Status:   status,
		Priority: priority,
		DateType: dateType,
	}
	sort := usecase.TaskSort{
		Field: sortField,
		Order: sortOrder,
	}
	return a.taskUseCase.GetFilteredAndSortedTasks(a.ctx, filter, sort)
}

func (a *App) GetTask(id string) (*domain.Task, error) {
	return a.taskUseCase.GetTask(a.ctx, id)
}

func (a *App) UpdateTask(id, title, description string) (*domain.Task, error) {
	return a.taskUseCase.UpdateTask(a.ctx, id, title, description)
}

func (a *App) ToggleTaskStatus(id string) (*domain.Task, error) {
	return a.taskUseCase.ToggleTaskStatus(a.ctx, id)
}

func (a *App) SetTaskPriority(id, priority string) (*domain.Task, error) {
	return a.taskUseCase.SetTaskPriority(a.ctx, id, priority)
}

func (a *App) SetTaskDueDate(id string, dueDate *time.Time) (*domain.Task, error) {
	return a.taskUseCase.SetTaskDueDate(a.ctx, id, dueDate)
}

func (a *App) DeleteTask(id string) error {
	return a.taskUseCase.DeleteTask(a.ctx, id)
}
