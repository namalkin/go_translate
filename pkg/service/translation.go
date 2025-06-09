package service

import (
	"encoding/json"
	"fmt"

	"github.com/namalkin/go_translate/pkg/repository"
	"github.com/namalkin/go_translate/pkg/tables"
)

type TranslationService struct {
	repo  repository.Translation
	cache repository.RedisCache
}

func NewTranslationService(repo repository.Translation, cache repository.RedisCache) *TranslationService {
	return &TranslationService{repo: repo, cache: cache}
}

func (s *TranslationService) Create(userId string, translation tables.Translation) (string, error) {
	return s.repo.Create(userId, translation)
}

func (s *TranslationService) GetAll(userId string, limit int) ([]tables.Translation, bool, error) {
	limitKey := "all"
	if limit > 0 {
		limitKey = fmt.Sprintf("%d", limit)
	}
	cacheKey := fmt.Sprintf("translations:all:%s:limit:%s", userId, limitKey)
	if s.cache != nil {
		if cached, err := s.cache.Get(cacheKey); err == nil && cached != "" {
			var result []tables.Translation
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				return result, true, nil
			}
		}
	}
	translations, err := s.repo.GetAll(userId)
	if err != nil {
		return nil, false, err
	}
	if limit > 0 && limit < len(translations) {
		translations = translations[:limit]
	}
	if s.cache != nil {
		_ = s.cache.Set(cacheKey, translations, 60) // 60 сек
	}
	return translations, false, err
}

func (s *TranslationService) GetById(userId, translationId string) (tables.Translation, bool, error) {
	cacheKey := fmt.Sprintf("translations:id:%s:%s", userId, translationId)
	if s.cache != nil {
		if cached, err := s.cache.Get(cacheKey); err == nil && cached != "" {
			var result tables.Translation
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				return result, true, nil
			}
		}
	}
	translation, err := s.repo.GetById(userId, translationId)
	if err == nil && s.cache != nil {
		_ = s.cache.Set(cacheKey, translation, 60)
	}
	return translation, false, err
}

func (s *TranslationService) Delete(userId, translationId string) error {
	if s.cache != nil {
		_ = s.cache.Del(fmt.Sprintf("translations:id:%s:%s", userId, translationId))
		_ = s.cache.Del(fmt.Sprintf("translations:all:%s", userId))
	}
	return s.repo.Delete(userId, translationId)
}

func (s *TranslationService) DeleteByPhrase(userId, phrase string) error {
	if s.cache != nil {
		_ = s.cache.Del(fmt.Sprintf("translations:all:%s", userId))
	}
	return s.repo.DeleteByPhrase(userId, phrase)
}

func (s *TranslationService) Update(userId, translationId string, input tables.UpdateTranslationInput) error {
	if s.cache != nil {
		_ = s.cache.Del(fmt.Sprintf("translations:id:%s:%s", userId, translationId))
		_ = s.cache.Del(fmt.Sprintf("translations:all:%s", userId))
	}
	return s.repo.Update(userId, translationId, input)
}
