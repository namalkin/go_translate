package service

import (
	"github.com/namalkin/go_translate/pkg/repository"
	"github.com/namalkin/go_translate/pkg/tables"
)

type Authorisation interface {
	CreateUser(user tables.User) (string, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (string, error)
}

type Translation interface {
	Create(userId string, translation tables.Translation) (string, error)
	GetAll(userId string, limit int) ([]tables.Translation, bool, error)
	GetById(userId, translationId string) (tables.Translation, error)
	Delete(userId, translationId string) error
	Update(userId, translationId string, input tables.UpdateTranslationInput) error
	DeleteByPhrase(userId, phrase string) error
}

type Service struct {
	Authorisation
	Translation
	Redis repository.RedisCache
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorisation: NewAuthService(repos.Authorisation),
		Translation:   NewTranslationService(repos.Translation, repos.Redis),
		Redis:         repos.Redis,
	}
}
