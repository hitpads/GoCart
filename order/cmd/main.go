package main

import (
	"GoCart/order/internal/adapters/http"
	"GoCart/order/internal/adapters/repository"
	"GoCart/order/internal/usecases"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := repository.NewInMemoryOrderRepository()
	orderUC := usecases.NewOrderUseCase(repo)

	r := gin.Default()
	http.RegisterOrderRoutes(r, orderUC)
	r.Run(":8002") // run on port 8002
}
