package domain

import (
	"time"
)

type Priority string

const (
	LowPriority    Priority = "low"
	MediumPriority Priority = "medium"
	HighPriority   Priority = "high"
)

type TaskStatus string

const (
	ActiveTask    TaskStatus = "active"
	CompletedTask TaskStatus = "completed"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewTask(title, description string) *Task {
	now := time.Now()
	return &Task{
		ID:          generateID(),
		Title:       title,
		Description: description,
		Status:      ActiveTask,
		Priority:    MediumPriority,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *Task) MarkComplete() {
	t.Status = CompletedTask
	t.UpdatedAt = time.Now()
}

func (t *Task) MarkActive() {
	t.Status = ActiveTask
	t.UpdatedAt = time.Now()
}

func (t *Task) SetPriority(priority Priority) {
	t.Priority = priority
	t.UpdatedAt = time.Now()
}

func (t *Task) SetDueDate(dueDate *time.Time) {
	t.DueDate = dueDate
	t.UpdatedAt = time.Now()
}

func (t *Task) IsOverdue() bool {
	if t.DueDate == nil || t.Status == CompletedTask {
		return false
	}
	return time.Now().After(*t.DueDate)
}