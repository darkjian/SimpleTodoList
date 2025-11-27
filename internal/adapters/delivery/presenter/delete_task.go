package presenter

import (
	"errors"
	"net/http"

	"github.com/darkjian/simpletodolist/internal/usecase"
	e "github.com/darkjian/simpletodolist/internal/usecase/errorx"
	"github.com/gin-gonic/gin"
)

func DeleteTaskWithError(c *gin.Context, err error) {
	var ex *e.ErrorX
	var statusCode = http.StatusInternalServerError
	var message = gin.H{"error": "internal server error"}

	if errors.As(err, &ex) {
		switch ex.Code {
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
