package handlers

import (
	"net/http"

	"metrics/internal/core/service"
	"metrics/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemHandler struct {
	service service.Pinger
}

func NewSystemHandler(service *service.SystemService) *SystemHandler {
	return &SystemHandler{service: service}
}

func (h *SystemHandler) Ping(ctx *gin.Context) {
	err := h.service.Ping(ctx)
	if err != nil {
		logger.Log.Error("Connectin to DB Error", zap.Error(err))
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	ctx.String(http.StatusOK, "OK")
}
