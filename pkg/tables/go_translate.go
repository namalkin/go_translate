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

type UpdateTranslationInput struct {
	Translation *string `json:"translation"`
	Done        *bool   `json:"done"`
}

func (i UpdateTranslationInput) Validate(expectedTranslation string) error {
	if i.Translation == nil {
		return errors.New("поле translation обязательно")
	}

	if *i.Translation != expectedTranslation {
		return errors.New("неверный перевод")
	}

	return nil
}
