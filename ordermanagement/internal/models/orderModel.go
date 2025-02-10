package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Order represents the order table in PostgreSQL
type Order struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Amount    float64   `gorm:"type:decimal(10,2);not null"`
	Status    string    `gorm:"type:varchar(20);not null;default:'pending'"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`

	// This defines the relationship with User
	User User `gorm:"foreignKey:user_id"`
}


// BeforeCreate runs before inserting an order
func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}
