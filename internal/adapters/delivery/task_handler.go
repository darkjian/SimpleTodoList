package delivery

import (
	"context"

	"github.com/darkjian/simpletodolist/internal/adapters/delivery/presenter"
	"github.com/darkjian/simpletodolist/internal/usecase"
	e "github.com/darkjian/simpletodolist/internal/usecase/errorx"
	"github.com/gin-gonic/gin"
)

type TaskServicer interface {
	GetTask(ctx context.Context, id string) (usecase.GetTaskOut, error)
	CreateTask(ctx context.Context, title string) (usecase.CreateTaskOut, error)
}

type TaskHandler struct {
	service TaskServicer
}

func NewTaskHandler(s TaskServicer) *TaskHandler {
	return &TaskHandler{service: s}
}

type GetTaskRequest struct {
	TaskID string `uri:"id" binding:"required,uuid"`
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	var request GetTaskRequest

	if err := c.ShouldBindUri(&request); err != nil {
		presenter.GetTaskWithError(c, e.New(usecase.CodeInvalidIDFormat, err))
		return
	}

	out, err := h.service.GetTask(c.Request.Context(), request.TaskID)
	if err != nil {
		presenter.GetTaskWithError(c, err)
		return
	}

	presenter.GetTaskOnSuccess(c, out)
}

type CreateTaskRequest struct {
	Title string `json:"title" binding:"required,min=1,max=20"`
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var request CreateTaskRequest

	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		presenter.CreateTaskWithError(c, e.New(usecase.CodeInvalidTitle, err))
		return
	}

	out, err := h.service.CreateTask(c.Request.Context(), request.Title)
	if err != nil {
		presenter.CreateTaskWithError(c, err)
		return
	}

	presenter.CreateTaskOnSuccess(c, out)
}
