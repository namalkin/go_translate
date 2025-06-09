package repository

import (
	"context"
	"errors"
	"time"

	"github.com/namalkin/go_translate/pkg/tables"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TranslationMongo struct {
	collection *mongo.Collection
}

func NewTranslationMongo(client *mongo.Client, dbName, collectionName string) *TranslationMongo {
	return &TranslationMongo{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

// создаёт новый перевод
func (r *TranslationMongo) Create(userId string, translation tables.Translation) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	translation.Done = false
	doc := bson.M{
		"phrase":               translation.Phrase,
		"expected_translation": translation.ExpectedTranslation,
		"done":                 translation.Done,
		"user_id":              userId,
	}

	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("не удалось преобразовать InsertedID в ObjectID")
	}
	return id.Hex(), nil
}

// возвращает все списки для пользователя
func (r *TranslationMongo) GetAll(userId string) ([]tables.Translation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var lists []tables.Translation
	if err = cursor.All(ctx, &lists); err != nil {
		return nil, err
	}
	return lists, nil
}

// возвращает список по его id и userId
func (r *TranslationMongo) GetById(userId, listId string) (tables.Translation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Преобразуем listId в ObjectID
	objectId, err := primitive.ObjectIDFromHex(listId)
	if err != nil {
		return tables.Translation{}, err
	}

	filter := bson.M{"_id": objectId, "user_id": userId}
	var list tables.Translation
	err = r.collection.FindOne(ctx, filter).Decode(&list)
	return list, err
}

// удаляет список по id
func (r *TranslationMongo) Delete(userId, listId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(listId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectId, "user_id": userId}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}

// обновляет список по id
func (r *TranslationMongo) Update(userId, listId string, input tables.UpdateTranslationInput) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// listId в ObjectID
	objectId, err := primitive.ObjectIDFromHex(listId)
	if err != nil {
		return err
	}

	// Сначала получаем текущий набор переводов
	var currentSet tables.Translation
	filter := bson.M{"_id": objectId, "user_id": userId}
	err = r.collection.FindOne(ctx, filter).Decode(&currentSet)
	if err != nil {
		return err
	}

	// Проверяем, совпадает ли перевод
	done := currentSet.ExpectedTranslation == *input.Translation

	// Обновляем только поле done
	update := bson.M{"$set": bson.M{"done": done}}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *TranslationMongo) DeleteByPhrase(userId, phrase string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": userId,
		"phrase":  phrase,
	}

	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}
