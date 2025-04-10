package repository

import (
	"GoCart/order/internal/domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface {
	Create(o *domain.Order) error
	GetByID(id string) (*domain.Order, error)
	UpdateStatus(id string, status domain.OrderStatus) error
	ListByUser(userID string) ([]domain.Order, error)
}

type MongoOrderRepository struct {
	collection *mongo.Collection
}

func NewMongoOrderRepository(client *mongo.Client, dbName string) *MongoOrderRepository {
	collection := client.Database(dbName).Collection("orders")
	return &MongoOrderRepository{
		collection: collection,
	}
}

func (repo *MongoOrderRepository) Create(o *domain.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.collection.InsertOne(ctx, o)
	return err
}

func (repo *MongoOrderRepository) GetByID(id string) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var order domain.Order
	filter := bson.M{"id": id}
	err := repo.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

func (repo *MongoOrderRepository) UpdateStatus(id string, status domain.OrderStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"status": status}}
	res, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (repo *MongoOrderRepository) ListByUser(userID string) ([]domain.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []domain.Order
	for cursor.Next(ctx) {
		var order domain.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
