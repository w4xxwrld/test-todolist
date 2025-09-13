package repository

import (
	"context"
	"time"

	"todo-list/internal/domain"
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id string) (*domain.Task, error)
	GetAll(ctx context.Context) ([]*domain.Task, error)
	GetByStatus(ctx context.Context, status domain.TaskStatus) ([]*domain.Task, error)
	GetByPriority(ctx context.Context, priority domain.Priority) ([]*domain.Task, error)
	GetByDateRange(ctx context.Context, from, to time.Time) ([]*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id string) error
}