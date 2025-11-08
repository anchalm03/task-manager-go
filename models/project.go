package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Id          uuid.UUID  `gorm:"column:id;default:gen_random_uuid();primaryKey"`
	Name        string     `gorm:"column:name;not null"`
	Description string     `gorm:"column:description"`
	OwnerID     *uuid.UUID `gorm:"column:owner_id"` // nullable due to ON DELETE SET NULL
}
