package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"todo-list/internal/domain"
)

type FileTaskRepository struct {
	filePath string
	tasks    map[string]*domain.Task
	mutex    sync.RWMutex
}

func NewFileTaskRepository(dataDir string) (*FileTaskRepository, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	filePath := filepath.Join(dataDir, "tasks.json")

	repo := &FileTaskRepository{
		filePath: filePath,
		tasks:    make(map[string]*domain.Task),
	}

	if err := repo.loadFromFile(); err != nil {
		return nil, fmt.Errorf("failed to load tasks from file: %w", err)
	}

	return repo, nil
}

func (r *FileTaskRepository) loadFromFile() error {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	var tasks []*domain.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return err
	}

	r.tasks = make(map[string]*domain.Task)
	for _, task := range tasks {
		r.tasks[task.ID] = task
	}

	return nil
}

func (r *FileTaskRepository) saveToFile() error {
	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	tempFile := r.filePath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return err
	}

	return os.Rename(tempFile, r.filePath)
}

func (r *FileTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.tasks[task.ID] = task
	return r.saveToFile()
}

func (r *FileTaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task with id %s not found", id)
	}

	return task, nil
}

func (r *FileTaskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})

	return tasks, nil
}

func (r *FileTaskRepository) GetByStatus(ctx context.Context, status domain.TaskStatus) ([]*domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tasks := make([]*domain.Task, 0)
	for _, task := range r.tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})

	return tasks, nil
}

func (r *FileTaskRepository) GetByPriority(ctx context.Context, priority domain.Priority) ([]*domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tasks := make([]*domain.Task, 0)
	for _, task := range r.tasks {
		if task.Priority == priority {
			tasks = append(tasks, task)
		}
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})

	return tasks, nil
}

func (r *FileTaskRepository) GetByDateRange(ctx context.Context, from, to time.Time) ([]*domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tasks := make([]*domain.Task, 0)
	for _, task := range r.tasks {
		if task.CreatedAt.After(from) && task.CreatedAt.Before(to) {
			tasks = append(tasks, task)
		}
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})

	return tasks, nil
}

func (r *FileTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return fmt.Errorf("task with id %s not found", task.ID)
	}

	r.tasks[task.ID] = task
	return r.saveToFile()
}

func (r *FileTaskRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return fmt.Errorf("task with id %s not found", id)
	}

	delete(r.tasks, id)
	return r.saveToFile()
}
