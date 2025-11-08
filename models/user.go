package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id       uuid.UUID `gorm:"column:id;default:gen_random_uuid();primaryKey"`
	Name     string    `gorm:"column:name;not null"`
	Email    string    `gorm:"column:email;not null;unique"`
	Password string    `gorm:"column:password;not null"`
	Role     string    `gorm:"column:role;default:'member'"`
}
