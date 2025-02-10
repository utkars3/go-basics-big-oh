package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name   string    `gorm:"type:varchar(100);not null"`
	Email  string    `gorm:"type:varchar(100);not null"`
	Mobile string    `gorm:"type:varchar(15);not null"`

	// Orders []Order `gorm:"foreignKey:UserID"`
}

func (u *User) ValidateUser() error {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			// You can add custom error handling here if needed
			fmt.Println(e.StructNamespace()) // Or more detailed output
		}
		return err
	}
	return nil
}
