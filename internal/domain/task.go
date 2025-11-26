package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrInternal         = errors.New("internal")
	ErrContextCancelled = errors.New("context cancelled")
	ErrAlreadyExists    = errors.New("already exists")
)

type ID string

type TaskRepo interface {
	ListTasks(ctx context.Context, limit, offset int) ([]Task, error)
	GetTask(ctx context.Context, id ID) (*Task, error)
	CreateTask(ctx context.Context, t *Task) error
	UpdateTask(ctx context.Context, t *Task) error
	SoftDeleteTask(ctx context.Context, id ID) error
}

type Task struct {
	ID          ID
	Title       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
	DeletedAt   *time.Time
}
