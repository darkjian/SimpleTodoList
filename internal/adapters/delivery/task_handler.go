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
	ListTasks(ctx context.Context, in usecase.ListTasksIn) (usecase.ListTasksOut, error)
	DeleteTask(ctx context.Context, id string) error
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

type ListTasksRequest struct {
	Limit  int `form:"limit" binding:"omitempty,gte=0,lte=100"`
	Offset int `form:"offset" binding:"omitempty,gte=0"`
}

func (r *ListTasksRequest) Process() {
	const defaultLimit = 10
	const defaultOffset = 0
	const maxLimit = 100

	if r.Limit <= 0 {
		r.Limit = defaultLimit
	}
	if r.Limit > maxLimit {
		r.Limit = maxLimit
	}
	if r.Offset < 0 {
		r.Offset = defaultOffset
	}
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	var request ListTasksRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		presenter.ListTasksWithError(c, e.New(usecase.CodeInvalidParam, err))
		return
	}

	request.Process()

	out, err := h.service.ListTasks(c.Request.Context(), usecase.ListTasksIn{
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		presenter.ListTasksWithError(c, err)
		return
	}

	presenter.ListTasksOnSuccess(c, out)
}

type DeleteTask struct {
	TaskID string `uri:"id" binding:"required,uuid"`
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	var request DeleteTask

	if err := c.ShouldBindUri(&request); err != nil {
		presenter.DeleteTaskWithError(c, e.New(usecase.CodeInvalidIDFormat, err))
		return
	}

	err := h.service.DeleteTask(c.Request.Context(), request.TaskID)
	if err != nil {
		presenter.DeleteTaskWithError(c, err)
		return
	}

	c.Status(200)
}
