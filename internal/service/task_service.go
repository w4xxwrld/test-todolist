package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"todo-list/internal/domain"
	"todo-list/internal/repository"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, title, description string) (*domain.Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("task title cannot be empty")
	}

	task := domain.NewTask(title, description)

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]*domain.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) GetTasksByStatus(ctx context.Context, status domain.TaskStatus) ([]*domain.Task, error) {
	return s.repo.GetByStatus(ctx, status)
}

func (s *TaskService) GetTasksByPriority(ctx context.Context, priority domain.Priority) ([]*domain.Task, error) {
	return s.repo.GetByPriority(ctx, priority)
}

func (s *TaskService) GetOverdueTasks(ctx context.Context) ([]*domain.Task, error) {
	allTasks, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	overdueTasks := make([]*domain.Task, 0)
	for _, task := range allTasks {
		if task.IsOverdue() {
			overdueTasks = append(overdueTasks, task)
		}
	}

	return overdueTasks, nil
}

func (s *TaskService) GetTodayTasks(ctx context.Context) ([]*domain.Task, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	allTasks, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	todayTasks := make([]*domain.Task, 0)
	for _, task := range allTasks {
		if task.DueDate != nil && task.DueDate.After(startOfDay) && task.DueDate.Before(endOfDay) {
			todayTasks = append(todayTasks, task)
		}
	}

	return todayTasks, nil
}

func (s *TaskService) GetWeekTasks(ctx context.Context) ([]*domain.Task, error) {
	now := time.Now()
	weekFromNow := now.Add(7 * 24 * time.Hour)

	allTasks, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	weekTasks := make([]*domain.Task, 0)
	for _, task := range allTasks {
		if task.DueDate != nil && task.DueDate.After(now) && task.DueDate.Before(weekFromNow) {
			weekTasks = append(weekTasks, task)
		}
	}

	return weekTasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id, title, description string) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(title) == "" {
		return nil, errors.New("task title cannot be empty")
	}

	task.Title = title
	task.Description = description
	task.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) MarkTaskComplete(ctx context.Context, id string) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.MarkComplete()

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) MarkTaskActive(ctx context.Context, id string) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.MarkActive()

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) SetTaskPriority(ctx context.Context, id string, priority domain.Priority) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.SetPriority(priority)

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) SetTaskDueDate(ctx context.Context, id string, dueDate *time.Time) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.SetDueDate(dueDate)

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}