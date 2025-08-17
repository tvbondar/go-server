package api

import (
    "net/http"
    "go-server/internal/service"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type Handler struct {
    service *service.OrderService
    logger  *zap.Logger
}

func NewHandler(service *service.OrderService, logger *zap.Logger) *Handler {
    return &Handler{service: service, logger: logger}
}

func (h *Handler) GetOrder(c *gin.Context) {
    orderUID := c.Param("order_uid")
    order, err := h.service.GetOrder(c.Request.Context(), orderUID)
    if err != nil {
        h.logger.Error("Failed to get order", zap.String("order_uid", orderUID), zap.Error(err))
        c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        return
    }
    c.JSON(http.StatusOK, order)
}

func (h *Handler) SetupRoutes(r *gin.Engine) {
    r.GET("/order/:order_uid", h.GetOrder)
    r.Static("/static", "./static")
    r.GET("/", func(c *gin.Context) {
        c.File("./static/index.html")
    })
}
