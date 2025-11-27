//go:generate mockgen -destination=./mocks/mock_task_repository.go -package=mocks github.com/darkjian/simpletodolist/internal/adapters/database TaskRepository
package database

import (
	"context"
	"time"

	"github.com/darkjian/simpletodolist/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository interface {
	domain.TaskRepo
}

type TaskRepo struct {
	db Database
}

type Database interface {
	Conn() *pgxpool.Pool
}

func NewTaskRepository(db Database) TaskRepository {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) ListTasks(ctx context.Context, limit int, offset int) (_ *domain.PaginatedTasks, _ error) {
	var total int

	err := r.db.Conn().QueryRow(ctx, `
		SELECT COUNT(*)
		FROM core.tasks
	`).Scan(&total)
	if err != nil {
		return nil, normalizePGError(err)
	}

	if offset >= total {
		return &domain.PaginatedTasks{Tasks: []domain.Task{}, Total: total}, nil
	}

	rows, err := r.db.Conn().Query(ctx, `
		SELECT id, title, created_at, updated_at, completed_at, deleted_at
		FROM core.tasks 
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		return nil, normalizePGError(err)
	}
	defer rows.Close()

	tasks := []domain.Task{}
	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.CompletedAt,
			&task.DeletedAt,
		); err != nil {
			return nil, normalizePGError(err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, normalizePGError(err)
	}

	nextOffset := offset + len(tasks)
	if nextOffset >= total {
		nextOffset = -1
	}

	return &domain.PaginatedTasks{
		Tasks:      tasks,
		Total:      total,
		NextOffset: nextOffset,
	}, nil
}

func (r *TaskRepo) GetTask(ctx context.Context, id domain.ID) (_ *domain.Task, _ error) {
	query := `SELECT 
	id, title, created_at, updated_at, completed_at, deleted_at 
	FROM core.tasks 
	WHERE id=$1`

	var task domain.Task

	if err := r.db.Conn().QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.CompletedAt,
		&task.DeletedAt,
	); err != nil {
		return nil, normalizePGError(err)
	}

	return &task, nil
}

func (r *TaskRepo) CreateTask(ctx context.Context, t *domain.Task) (_ error) {
	query := `INSERT 
	INTO core.tasks (title)
	VALUES ($1)
	RETURNING id, created_at, updated_at;`

	if err := r.db.Conn().QueryRow(ctx, query, t.Title).Scan(
		&t.ID,
		&t.CreatedAt,
		&t.UpdatedAt,
	); err != nil {
		return normalizePGError(err)
	}

	return nil
}

func (r *TaskRepo) CompleteTask(ctx context.Context, id domain.ID, completedAt time.Time) (_ error) {
	query := `UPDATE core.tasks SET completed_at=COALESCE(completed_at, $1) WHERE id=$2`
	if _, err := r.db.Conn().Exec(ctx, query, completedAt, id); err != nil {
		return normalizePGError(err)
	}

	return nil
}

func (r *TaskRepo) DeleteTask(ctx context.Context, id domain.ID) (_ error) {
	query := `DELETE from core.tasks WHERE id=$1`
	if _, err := r.db.Conn().Exec(ctx, query, id); err != nil {
		return normalizePGError(err)
	}

	return nil
}
