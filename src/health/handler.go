package health

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ctx context.Context
}

func NewHandler(ctx context.Context) *Handler {
	return &Handler{
		ctx: ctx,
	}
}

func (handler *Handler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "running"})
}
