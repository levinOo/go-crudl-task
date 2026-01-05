package handlers

import (
	"github.com/levinOo/go-crudl-task/internal/domain"

	"github.com/gin-gonic/gin"
)

// newErrorResponse отправляет ошибку в JSON и логирует её (опционально)
func newErrorResponse(c *gin.Context, statusCode int, message string) {

	c.AbortWithStatusJSON(statusCode, domain.ErrorResponse{
		Error: message,
	})
}
