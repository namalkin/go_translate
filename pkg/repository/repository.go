package repository

import (
	"github.com/namalkin/go_translate/pkg/tables"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorisation interface {
	CreateUser(user tables.User) (string, error)
	GetUser(username, password string) (tables.User, error)
}

type Translation interface {
	Create(userId string, translation tables.Translation) (string, error)
	GetAll(userId string) ([]tables.Translation, error)
	GetById(userId, translationId string) (tables.Translation, error)
	Delete(userId, translationId string) error
	Update(userId, translationId string, input tables.UpdateTranslationInput) error
	DeleteByPhrase(userId, phrase string) error
}

type RedisCache interface {
	Set(key string, value interface{}, ttlSeconds int) error
	Get(key string) (string, error)
	Del(key string) error
}

type Repository struct {
	Authorisation
	Translation
	Redis RedisCache
}

func NewRepository(db *mongo.Client, dbName, collectionName string, redis RedisCache) *Repository {
	return &Repository{
		Authorisation: NewAuthMongo(db, dbName, collectionName),
		Translation:   NewTranslationMongo(db, dbName, "translations"),
		Redis:         redis,
	}
}
