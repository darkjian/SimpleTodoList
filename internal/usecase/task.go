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
