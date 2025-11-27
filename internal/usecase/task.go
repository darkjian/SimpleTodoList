package usecase

import (
	"context"
	"errors"

	"github.com/darkjian/simpletodolist/internal/domain"
	e "github.com/darkjian/simpletodolist/internal/usecase/errorx"
	"github.com/google/uuid"
)

type GetTaskRepository interface {
	GetTask(ctx context.Context, id domain.ID) (_ *domain.Task, _ error)
	CreateTask(ctx context.Context, t *domain.Task) error
	ListTasks(ctx context.Context, limit, offset int) (*domain.PaginatedTasks, error)
	SoftDeleteTask(ctx context.Context, id domain.ID) error
}

type TaskService struct {
	repo GetTaskRepository
}

func NewTaskService(r GetTaskRepository) *TaskService {
	return &TaskService{repo: r}
}

type GetTaskOut struct {
	Task domain.Task
}

func (s *TaskService) GetTask(ctx context.Context, id string) (GetTaskOut, error) {
	if _, err := uuid.Parse(id); err != nil {
		return GetTaskOut{}, e.New(CodeInvalidIDFormat, errors.New("invalid UUID format"))
	}

	task, err := s.repo.GetTask(ctx, domain.ID(id))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return GetTaskOut{}, e.New(CodeNotFound, err)
		default:
			return GetTaskOut{}, e.New(CodeInternal, err)
		}
	}

	return GetTaskOut{Task: *task}, nil
}

type CreateTaskOut struct {
	Task domain.Task
}

func (s *TaskService) CreateTask(ctx context.Context, title string) (CreateTaskOut, error) {

	titleLen := len(title)
	if titleLen == 0 || titleLen > 20 {
		return CreateTaskOut{}, e.New(CodeInvalidTitle, errors.New("incorrect title length"))
	}

	task := &domain.Task{
		Title: title,
	}

	if err := s.repo.CreateTask(ctx, task); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return CreateTaskOut{}, e.New(CodeNotFound, err)
		default:
			return CreateTaskOut{}, e.New(CodeInternal, err)
		}
	}

	return CreateTaskOut{Task: *task}, nil
}

type ListTasksIn struct {
	Limit, Offset int
}

type ListTasksOut struct {
	Tasks             []domain.Task
	Total, NextOffset int
}

func (s *TaskService) ListTasks(ctx context.Context, in ListTasksIn) (ListTasksOut, error) {

	out, err := s.repo.ListTasks(ctx, in.Limit, in.Offset)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return ListTasksOut{}, e.New(CodeNotFound, err)
		default:
			return ListTasksOut{}, e.New(CodeInternal, err)
		}
	}

	return ListTasksOut{
		Tasks:      out.Tasks,
		Total:      out.Total,
		NextOffset: out.NextOffset,
	}, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return e.New(CodeInvalidIDFormat, errors.New("invalid UUID format"))
	}

	if err := s.repo.SoftDeleteTask(ctx, domain.ID(id)); err != nil {
		return e.New(CodeInternal, err)
	}

	return nil
}
