package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Inventory struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name  string    `gorm:"type:varchar(100);not null" validate:"required"`
	Stock int       `gorm:"type:int;not null" validate:"required,min=0"`
}

func (i *Inventory) ValidateInventory() error {
	validate := validator.New()
	err := validate.Struct(i)
	if err != nil {
		return err
	}
	return nil
}
