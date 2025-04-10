package http

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"GoCart/inventory/internal/domain"
	"GoCart/inventory/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productUC usecases.ProductUseCase
}

// check product related routes
func RegisterProductRoutes(router *gin.Engine, productUC usecases.ProductUseCase) {
	handler := &ProductHandler{productUC: productUC}

	v1 := router.Group("/products")
	{
		v1.POST("", handler.CreateProduct)
		v1.GET("/:id", handler.GetProduct)
		v1.PATCH("/:id", handler.UpdateProduct)
		v1.DELETE("/:id", handler.DeleteProduct)
		v1.GET("", handler.ListProducts)
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var p domain.Product
	// read / log raw body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	log.Printf("Raw request body: %s", string(bodyBytes))

	// reset body so it can be read again by ShouldBindJSON
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = uuid.New().String()
	log.Printf("Creating product: %+v", p)

	if err := h.productUC.CreateProduct(&p); err != nil {
		log.Printf("Failed to create product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	p, err := h.productUC.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var p domain.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = id
	if err := h.productUC.UpdateProduct(id, &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.productUC.DeleteProduct(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.productUC.ListProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}
