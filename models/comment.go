package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Id     uuid.UUID  `gorm:"column:id;default:gen_random_uuid();primaryKey"`
	Text   string     `gorm:"column:text;not null"`
	TaskID uuid.UUID  `gorm:"column:task_id"` // CASCADE delete
	UserID *uuid.UUID `gorm:"column:user_id"` // nullable due to ON DELETE SET NULL
}
