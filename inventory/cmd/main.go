package main

import (
	"GoCart/inventory/internal/adapters/http"
	"GoCart/inventory/internal/adapters/repository"
	"GoCart/inventory/internal/usecases"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := repository.NewInMemoryProductRepository()
	productUC := usecases.NewProductUseCase(repo)

	r := gin.Default()
	http.RegisterProductRoutes(r, productUC)

	r.Run(":8001")
}
