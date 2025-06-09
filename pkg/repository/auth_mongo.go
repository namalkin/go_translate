package repository

import (
	"context"
	"time"

	"github.com/namalkin/go_translate/pkg/tables"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthMongo struct {
	collection *mongo.Collection
}

func NewAuthMongo(client *mongo.Client, dbName, collectionName string) *AuthMongo {
	collection := client.Database(dbName).Collection(collectionName)
	return &AuthMongo{collection: collection}
}

func (r *AuthMongo) CreateUser(user tables.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex() // строковое представление ObjectId

	return id, nil
}

func (r *AuthMongo) GetUser(username, password string) (tables.User, error) {
	var user tables.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"username": username, "password": password}
	err := r.collection.FindOne(ctx, filter).Decode(&user)

	return user, err
}
