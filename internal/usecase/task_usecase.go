package usecase

import (
	"context"
	"sort"
	"time"

	"todo-list/internal/domain"
	"todo-list/internal/service"
)

type TaskUseCase struct {
	taskService *service.TaskService
}

func NewTaskUseCase(taskService *service.TaskService) *TaskUseCase {
	return &TaskUseCase{
		taskService: taskService,
	}
}

type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type TaskFilter struct {
	Status   string `json:"status"`
	Priority string `json:"priority"`
	DateType string `json:"date"`
}

type TaskSort struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

func (uc *TaskUseCase) CreateTask(ctx context.Context, req CreateTaskRequest) (*domain.Task, error) {
	task, err := uc.taskService.CreateTask(ctx, req.Title, req.Description)
	if err != nil {
		return nil, err
	}

	if req.Priority != "" {
		priority := domain.Priority(req.Priority)
		if priority == domain.LowPriority || priority == domain.MediumPriority || priority == domain.HighPriority {
			task, err = uc.taskService.SetTaskPriority(ctx, task.ID, priority)
			if err != nil {
				return nil, err
			}
		}
	}

	if req.DueDate != nil {
		task, err = uc.taskService.SetTaskDueDate(ctx, task.ID, req.DueDate)
		if err != nil {
			return nil, err
		}
	}

	return task, nil
}

func (uc *TaskUseCase) GetFilteredAndSortedTasks(ctx context.Context, filter TaskFilter, sort TaskSort) ([]*domain.Task, error) {
	var tasks []*domain.Task
	var err error

	switch filter.Status {
	case "active":
		tasks, err = uc.taskService.GetTasksByStatus(ctx, domain.ActiveTask)
	case "completed":
		tasks, err = uc.taskService.GetTasksByStatus(ctx, domain.CompletedTask)
	default:
		tasks, err = uc.taskService.GetAllTasks(ctx)
	}

	if err != nil {
		return nil, err
	}

	// Filter by priority
	if filter.Priority != "all" && filter.Priority != "" {
		priorityFiltered := make([]*domain.Task, 0)
		priority := domain.Priority(filter.Priority)
		for _, task := range tasks {
			if task.Priority == priority {
				priorityFiltered = append(priorityFiltered, task)
			}
		}
		tasks = priorityFiltered
	}

	// Filter by date
	switch filter.DateType {
	case "today":
		todayTasks, err := uc.taskService.GetTodayTasks(ctx)
		if err != nil {
			return nil, err
		}
		tasks = uc.intersectTasks(tasks, todayTasks)
	case "week":
		weekTasks, err := uc.taskService.GetWeekTasks(ctx)
		if err != nil {
			return nil, err
		}
		tasks = uc.intersectTasks(tasks, weekTasks)
	case "overdue":
		overdueTasks, err := uc.taskService.GetOverdueTasks(ctx)
		if err != nil {
			return nil, err
		}
		tasks = uc.intersectTasks(tasks, overdueTasks)
	}

	// Apply sorting
	uc.sortTasks(tasks, sort)

	return tasks, nil
}

func (uc *TaskUseCase) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	return uc.taskService.GetTaskByID(ctx, id)
}

func (uc *TaskUseCase) UpdateTask(ctx context.Context, id, title, description string) (*domain.Task, error) {
	return uc.taskService.UpdateTask(ctx, id, title, description)
}

func (uc *TaskUseCase) ToggleTaskStatus(ctx context.Context, id string) (*domain.Task, error) {
	task, err := uc.taskService.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if task.Status == domain.ActiveTask {
		return uc.taskService.MarkTaskComplete(ctx, id)
	} else {
		return uc.taskService.MarkTaskActive(ctx, id)
	}
}

func (uc *TaskUseCase) SetTaskPriority(ctx context.Context, id, priority string) (*domain.Task, error) {
	return uc.taskService.SetTaskPriority(ctx, id, domain.Priority(priority))
}

func (uc *TaskUseCase) SetTaskDueDate(ctx context.Context, id string, dueDate *time.Time) (*domain.Task, error) {
	return uc.taskService.SetTaskDueDate(ctx, id, dueDate)
}

func (uc *TaskUseCase) DeleteTask(ctx context.Context, id string) error {
	return uc.taskService.DeleteTask(ctx, id)
}

func (uc *TaskUseCase) intersectTasks(tasks1, tasks2 []*domain.Task) []*domain.Task {
	taskMap := make(map[string]*domain.Task)
	for _, task := range tasks1 {
		taskMap[task.ID] = task
	}

	result := make([]*domain.Task, 0)
	for _, task := range tasks2 {
		if _, exists := taskMap[task.ID]; exists {
			result = append(result, task)
		}
	}

	return result
}

func (uc *TaskUseCase) sortTasks(tasks []*domain.Task, sortBy TaskSort) {
	switch sortBy.Field {
	case "priority":
		sort.Slice(tasks, func(i, j int) bool {
			priorityOrder := map[domain.Priority]int{
				domain.HighPriority:   3,
				domain.MediumPriority: 2,
				domain.LowPriority:    1,
			}

			if sortBy.Order == "desc" {
				return priorityOrder[tasks[i].Priority] > priorityOrder[tasks[j].Priority]
			}
			return priorityOrder[tasks[i].Priority] < priorityOrder[tasks[j].Priority]
		})
	case "due_date":
		sort.Slice(tasks, func(i, j int) bool {
			// Tasks without due dates go to the end
			if tasks[i].DueDate == nil && tasks[j].DueDate == nil {
				return false
			}
			if tasks[i].DueDate == nil {
				return false
			}
			if tasks[j].DueDate == nil {
				return true
			}

			if sortBy.Order == "desc" {
				return tasks[i].DueDate.After(*tasks[j].DueDate)
			}
			return tasks[i].DueDate.Before(*tasks[j].DueDate)
		})
	default: // "created"
		sort.Slice(tasks, func(i, j int) bool {
			if sortBy.Order == "desc" {
				return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
			}
			return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		})
	}
}
