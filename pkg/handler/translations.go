package handler

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/namalkin/go_translate/pkg/tables"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// POST новый перевод
func (h *Handler) createTranslation(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input tables.Translation
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Некорректный JSON")
		return
	}

	if input.Phrase == "" || input.ExpectedTranslation == "" {
		newErrorResponse(c, http.StatusBadRequest, "Фраза и ожидаемый перевод обязательны")
		return
	}

	id, err := h.services.Translation.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка создания перевода")
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// GET все переводы пользователя
func (h *Handler) getAllTranslations(c *gin.Context) {
	start := time.Now()
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	limit := 0
	limitStr := c.Query("limit")
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	translations, fromCache, err := h.services.Translation.GetAll(userId, limit)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка получения переводов")
		return
	}

	total := len(translations)
	source := "db"
	if fromCache {
		source = "cache"
	}

	duration := time.Since(start).Milliseconds()

	// консоль GIN
	c.Writer.WriteString("[GIN-debug] translations source: " + source + "\n")

	c.JSON(http.StatusOK, gin.H{
		"data":        translations,
		"total":       total,
		"source":      source,
		"duration_ms": duration,
	})
}

// GET перевод по id
func (h *Handler) getTranslationById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	translationId := c.Param("id")
	if translationId == "" {
		newErrorResponse(c, http.StatusBadRequest, "Пустой id параметр")
		return
	}

	if _, err := primitive.ObjectIDFromHex(translationId); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Некорректный id параметр")
		return
	}

	translation, err := h.services.Translation.GetById(userId, translationId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			newErrorResponse(c, http.StatusNotFound, "Перевод не найден")
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "Ошибка получения перевода")
		}
		return
	}

	c.JSON(http.StatusOK, translation)
}

// PUT перевод по id
func (h *Handler) updateTranslation(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	translationId := c.Param("id")
	if translationId == "" {
		newErrorResponse(c, http.StatusBadRequest, "Пустой id параметр")
		return
	}

	if _, err := primitive.ObjectIDFromHex(translationId); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Некорректный id параметр")
		return
	}

	var input tables.UpdateTranslationInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Некорректный JSON")
		return
	}

	if err := input.Validate(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Translation.Update(userId, translationId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка обновления перевода")
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

// DELETE перевод по id
func (h *Handler) deleteTranslation(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	translationId := c.Param("id")
	if translationId == "" {
		newErrorResponse(c, http.StatusBadRequest, "Пустой id параметр")
		return
	}

	if _, err := primitive.ObjectIDFromHex(translationId); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Некорректный id параметр")
		return
	}

	if err := h.services.Translation.Delete(userId, translationId); err != nil {
		if err == mongo.ErrNoDocuments {
			newErrorResponse(c, http.StatusNotFound, "Перевод не найден")
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "Ошибка удаления перевода")
		}
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

func (h *Handler) createTestTranslations(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	testData := []tables.Translation{
		{Phrase: "hello", ExpectedTranslation: "привет"},
		{Phrase: "world", ExpectedTranslation: "мир"},
		{Phrase: "sun", ExpectedTranslation: "солнце"},
		{Phrase: "moon", ExpectedTranslation: "луна"},
		{Phrase: "sky", ExpectedTranslation: "небо"},
		{Phrase: "tree", ExpectedTranslation: "дерево"},
		{Phrase: "flower", ExpectedTranslation: "цветок"},
		{Phrase: "water", ExpectedTranslation: "вода"},
		{Phrase: "fire", ExpectedTranslation: "огонь"},
		{Phrase: "earth", ExpectedTranslation: "земля"},
	}

	var createdIds []string
	baseWords := testData

	for i := 0; i < 10; i++ {
		for _, word := range baseWords {
			translation := tables.Translation{
				Phrase:              word.Phrase + "_" + string(rune('0'+i)),
				ExpectedTranslation: word.ExpectedTranslation + "_" + string(rune('0'+i)),
				Done:                false,
			}

			id, err := h.services.Translation.Create(userId, translation)
			if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, "Error creating test translation")
				return
			}
			createdIds = append(createdIds, id)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"count":  len(createdIds),
		"ids":    createdIds,
	})
}

// Тест: 100 000 переводов через API порционно и параллельно
func (h *Handler) createTestTranslations100k(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	baseWords := []tables.Translation{
		{Phrase: "hello", ExpectedTranslation: "привет"},
		{Phrase: "world", ExpectedTranslation: "мир"},
		{Phrase: "sun", ExpectedTranslation: "солнце"},
		{Phrase: "moon", ExpectedTranslation: "луна"},
		{Phrase: "sky", ExpectedTranslation: "небо"},
		{Phrase: "tree", ExpectedTranslation: "дерево"},
		{Phrase: "flower", ExpectedTranslation: "цветок"},
		{Phrase: "water", ExpectedTranslation: "вода"},
		{Phrase: "fire", ExpectedTranslation: "огонь"},
		{Phrase: "earth", ExpectedTranslation: "земля"},
	}

	const total = 100000
	const batchSize = 1000
	batchCount := total / batchSize

	var wg sync.WaitGroup
	var mu sync.Mutex
	created := 0
	errors := 0

	for b := 0; b < batchCount; b++ {
		wg.Add(1)
		go func(batchNum int) {
			defer wg.Done()
			localCreated := 0
			for i := 0; i < batchSize/len(baseWords); i++ {
				for j, word := range baseWords {
					idx := batchNum*batchSize + i*len(baseWords) + j
					translation := tables.Translation{
						Phrase:              word.Phrase + "_" + strconv.Itoa(idx),
						ExpectedTranslation: word.ExpectedTranslation + "_" + strconv.Itoa(idx),
						Done:                false,
					}
					_, err := h.services.Translation.Create(userId, translation)
					if err == nil {
						localCreated++
					} else {
						mu.Lock()
						errors++
						mu.Unlock()
					}
				}
			}
			mu.Lock()
			created += localCreated
			mu.Unlock()
		}(b)
	}

	wg.Wait()

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"count":  created,
		"errors": errors,
	})
}

// Тест: удалить все переводы пользователя порционно и параллельно
func (h *Handler) deleteAllTranslations(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	translations, _, err := h.services.Translation.GetAll(userId, 0)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка получения переводов")
		return
	}

	const batchSize = 1000
	var wg sync.WaitGroup
	var mu sync.Mutex
	deleted := 0

	for i := 0; i < len(translations); i += batchSize {
		end := i + batchSize
		if end > len(translations) {
			end = len(translations)
		}
		batch := translations[i:end]
		wg.Add(1)
		go func(batch []tables.Translation) {
			defer wg.Done()
			localDeleted := 0
			for _, t := range batch {
				_ = h.services.Translation.Delete(userId, t.Id.Hex())
				localDeleted++
			}
			mu.Lock()
			deleted += localDeleted
			mu.Unlock()
		}(batch)
	}

	wg.Wait()
	// TODO: Удалить все кеши translations:all:{userId}:limit:*
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"deleted": deleted,
	})
}
