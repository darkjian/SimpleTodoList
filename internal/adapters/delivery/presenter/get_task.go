package presenter

import (
	"errors"
	"net/http"

	"github.com/darkjian/simpletodolist/internal/usecase"
	e "github.com/darkjian/simpletodolist/internal/usecase/errorx"
	"github.com/gin-gonic/gin"
)

func GetTaskWithError(c *gin.Context, err error) {
	var ex *e.ErrorX
	var statusCode = http.StatusInternalServerError
	var message = gin.H{"error": "internal server error"}

	if errors.As(err, &ex) {
		switch ex.Code {
		case usecase.CodeNotFound:
			message["error"] = "not found"
			statusCode = 404
		case usecase.CodeInternal:
			message["error"] = "internal server error"
			statusCode = 500
		case usecase.CodeInvalidIDFormat:
			message["error"] = "invalid id"
			statusCode = 400
		}
	}

	c.AbortWithStatusJSON(statusCode, message)
}

func GetTaskOnSuccess(c *gin.Context, out usecase.GetTaskOut) {
	c.JSON(200, gin.H{
		"id":           out.Task.ID,
		"title":        out.Task.Title,
		"created_at":   out.Task.CreatedAt,
		"updated_at":   out.Task.UpdatedAt,
		"completed_at": out.Task.CompletedAt,
	})
}
