package repository

import (
	"GoCart/inventory/internal/domain"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	Create(p *domain.Product) error
	GetByID(id string) (*domain.Product, error)
	Update(id string, p *domain.Product) error
	Delete(id string) error
	List() ([]domain.Product, error)
}

type MongoProductRepository struct {
	collection *mongo.Collection
}

func NewMongoProductRepository(client *mongo.Client, dbName string) *MongoProductRepository {
	collection := client.Database(dbName).Collection("products")
	return &MongoProductRepository{
		collection: collection,
	}
}

func (repo *MongoProductRepository) Create(p *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Inserting product into MongoDB: %+v", p)
	result, err := repo.collection.InsertOne(ctx, p)
	if err != nil {
		log.Printf("Error inserting product: %v", err)
		return err
	}
	log.Printf("Product inserted with ID: %v", result.InsertedID)
	return nil
}

func (repo *MongoProductRepository) GetByID(id string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product domain.Product
	filter := bson.M{"id": id}
	err := repo.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (repo *MongoProductRepository) Update(id string, p *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": p}
	res, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (repo *MongoProductRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	res, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (repo *MongoProductRepository) List() ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	for cursor.Next(ctx) {
		var prod domain.Product
		if err := cursor.Decode(&prod); err != nil {
			return nil, err
		}
		products = append(products, prod)
	}

	return products, nil
}
