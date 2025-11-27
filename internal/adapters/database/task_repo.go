//go:generate mockgen -destination=./mocks/mock_task_repository.go -package=mocks github.com/darkjian/simpletodolist/internal/adapters/database TaskRepository
package database

import (
	"context"

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

func (r *TaskRepo) ListTasks(ctx context.Context, limit int, offset int) (_ []domain.Task, _ error) {
	panic("not implemented") // TODO: Implement
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

func (r *TaskRepo) UpdateTask(ctx context.Context, t *domain.Task) (_ error) {
	panic("not implemented") // TODO: Implement
}

func (r *TaskRepo) SoftDeleteTask(ctx context.Context, id domain.ID) (_ error) {
	panic("not implemented") // TODO: Implement
}
