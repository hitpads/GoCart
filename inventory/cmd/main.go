package main

import (
	"context"
	"log"
	"time"

	"GoCart/inventory/internal/adapters/http"
	"GoCart/inventory/internal/adapters/repository"
	"GoCart/inventory/internal/usecases"

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
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	repoMongo := repository.NewMongoProductRepository(client, "gocart")
	productUC := usecases.NewProductUseCase(repoMongo)

	r := gin.Default()
	http.RegisterProductRoutes(r, productUC)

	r.Run(":8001")
}
