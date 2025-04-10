package http

import (
	"net/http"

	"GoCart/order/internal/domain"
	"GoCart/order/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderUC usecases.OrderUseCase
}

func RegisterOrderRoutes(router *gin.Engine, orderUC usecases.OrderUseCase) {
	handler := &OrderHandler{orderUC: orderUC}

	v1 := router.Group("/orders")
	{
		v1.POST("", handler.CreateOrder)
		v1.GET("/:id", handler.GetOrder)
		v1.PATCH("/:id", handler.UpdateOrder)
		v1.GET("", handler.ListOrders)
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var o domain.Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// generate new order id
	o.ID = uuid.New().String()
	if err := h.orderUC.CreateOrder(&o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, o)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderUC.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status domain.OrderStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.orderUC.UpdateOrderStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}
	orders, err := h.orderUC.ListOrdersByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}
