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
