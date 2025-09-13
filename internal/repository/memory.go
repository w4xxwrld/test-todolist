package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"todo-list/internal/domain"
)

type MemoryTaskRepository struct {
	tasks map[string]*domain.Task
	mutex sync.RWMutex
}

func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

func (r *MemoryTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.tasks[task.ID] = task
	return nil
}

func (r *MemoryTaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task with id %s not found", id)
	}

	return task, nil
}

func (r *MemoryTaskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	// Sort by creation date (newest first)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})

	return tasks, nil
}

func (r *MemoryTaskRepository) GetByStatus(ctx context.Context, status domain.TaskStatus) ([]*domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tasks := make([]*domain.Task, 0)
	for _, task := range r.tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}

	// Sort by creation date (newest first)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})

	return tasks, nil
}

func (r *MemoryTaskRepository) GetByPriority(ctx context.Context, priority domain.Priority) ([]*domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tasks := make([]*domain.Task, 0)
	for _, task := range r.tasks {
		if task.Priority == priority {
			tasks = append(tasks, task)
		}
	}

	// Sort by creation date (newest first)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})

	return tasks, nil
}

func (r *MemoryTaskRepository) GetByDateRange(ctx context.Context, from, to time.Time) ([]*domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tasks := make([]*domain.Task, 0)
	for _, task := range r.tasks {
		if task.CreatedAt.After(from) && task.CreatedAt.Before(to) {
			tasks = append(tasks, task)
		}
	}

	// Sort by creation date (newest first)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})

	return tasks, nil
}

func (r *MemoryTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return fmt.Errorf("task with id %s not found", task.ID)
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *MemoryTaskRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return fmt.Errorf("task with id %s not found", id)
	}

	delete(r.tasks, id)
	return nil
}