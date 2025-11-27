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
	ListTasks(ctx context.Context, limit, offset int) (*PaginatedTasks, error)
	GetTask(ctx context.Context, id ID) (*Task, error)
	CreateTask(ctx context.Context, t *Task) error
	UpdateTask(ctx context.Context, t *Task) error
	SoftDeleteTask(ctx context.Context, id ID) error
}

type Task struct {
	ID          ID         `json:"id"`
	Title       string     `json:"title"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type PaginatedTasks struct {
	Tasks             []Task
	Total, NextOffset int
}
