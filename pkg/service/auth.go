package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/namalkin/go_translate/pkg/repository"
	"github.com/namalkin/go_translate/pkg/tables"
)

const (
	salt       = "erijj4or-3j4or34r"
	signingKey = "fjijerwqdqw-dewfewf"
	TokenTTL   = 12 * time.Hour
)

type AuthService struct {
	repo repository.Authorisation
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

func NewAuthService(repo repository.Authorisation) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user tables.User) (string, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID.Hex(),
	})

	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}
	return "Bearer " + signedToken, nil
}

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	t, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("ошибка signin вызова")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return "0", err
	}

	claims, ok := t.Claims.(*tokenClaims)
	if !ok {
		return "0", errors.New("токен не является типом *tokenClaims")
	}
	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New() // хеширование
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt))) // хэш+соль пароля
}
