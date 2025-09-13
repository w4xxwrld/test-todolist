package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"todo-list/internal/domain"

	_ "github.com/lib/pq"
)

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(connectionString string) (*PostgresTaskRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := &PostgresTaskRepository{
		db: db,
	}

	// Create tables if they don't exist
	if err := repo.createTables(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return repo, nil
}

func (r *PostgresTaskRepository) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id VARCHAR(255) PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT DEFAULT '',
		status VARCHAR(50) NOT NULL CHECK (status IN ('active', 'completed')),
		priority VARCHAR(50) NOT NULL CHECK (priority IN ('low', 'medium', 'high')),
		due_date TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
	CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);
	CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
	`

	_, err := r.db.Exec(query)
	return err
}

func (r *PostgresTaskRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, status, priority, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(
		ctx, query,
		task.ID,
		task.Title,
		task.Description,
		string(task.Status),
		string(task.Priority),
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	)

	return err
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var task domain.Task
	var status, priority string
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id %s not found", id)
		}
		return nil, err
	}

	task.Status = domain.TaskStatus(status)
	task.Priority = domain.Priority(priority)

	return &task, nil
}

func (r *PostgresTaskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`

	return r.queryTasks(ctx, query)
}

func (r *PostgresTaskRepository) GetByStatus(ctx context.Context, status domain.TaskStatus) ([]*domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE status = $1
		ORDER BY created_at DESC
	`

	return r.queryTasks(ctx, query, string(status))
}

func (r *PostgresTaskRepository) GetByPriority(ctx context.Context, priority domain.Priority) ([]*domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE priority = $1
		ORDER BY created_at DESC
	`

	return r.queryTasks(ctx, query, string(priority))
}

func (r *PostgresTaskRepository) GetByDateRange(ctx context.Context, from, to time.Time) ([]*domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE created_at BETWEEN $1 AND $2
		ORDER BY created_at DESC
	`

	return r.queryTasks(ctx, query, from, to)
}

func (r *PostgresTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	query := `
		UPDATE tasks
		SET title = $2, description = $3, status = $4, priority = $5, due_date = $6, updated_at = $7
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx, query,
		task.ID,
		task.Title,
		task.Description,
		string(task.Status),
		string(task.Priority),
		task.DueDate,
		task.UpdatedAt,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %s not found", task.ID)
	}

	return nil
}

func (r *PostgresTaskRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %s not found", id)
	}

	return nil
}

func (r *PostgresTaskRepository) queryTasks(ctx context.Context, query string, args ...interface{}) ([]*domain.Task, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task

	for rows.Next() {
		var task domain.Task
		var status, priority string

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&status,
			&priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		task.Status = domain.TaskStatus(status)
		task.Priority = domain.Priority(priority)

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}