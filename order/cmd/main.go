package main

import (
	"context"
	"log"
	"time"

	"GoCart/order/internal/adapters/http"
	"GoCart/order/internal/adapters/repository"
	"GoCart/order/internal/usecases"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	// verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	repoMongo := repository.NewMongoOrderRepository(client, "gocart")
	orderUC := usecases.NewOrderUseCase(repoMongo)

	r := gin.Default()
	http.RegisterOrderRoutes(r, orderUC)

	r.Run(":8002")
}
