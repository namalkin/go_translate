package tables

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Translation struct {
	Id                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Phrase              string             `json:"phrase" bson:"phrase" binding:"required"`
	ExpectedTranslation string             `json:"expected_translation" bson:"expected_translation" binding:"required"`
	Done                bool               `json:"done" bson:"done"`
}

type UserProgress struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	UserId        primitive.ObjectID `bson:"user_id"`
	TranslationId primitive.ObjectID `bson:"translation_id"`
	Done          bool               `bson:"done"`
}

type UpdateTranslationInput struct {
	Translation *string `json:"translation" binding:"required"`
}

func (i UpdateTranslationInput) Validate() error {
	if i.Translation == nil {
		return errors.New("translation is required")
	}
	return nil
}
