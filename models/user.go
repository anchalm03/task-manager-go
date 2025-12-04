package models

import (
	"task_manager/db"

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

func CreateUser(user *User) error {
	return db.DB.Create(user).Error
}

// Get user by email (for login)
func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := db.DB.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Get user by ID
func GetUserById(userID uuid.UUID) (*User, error) {
	var user User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Delete user
func DeleteUser(userID uuid.UUID) error {
	return db.DB.Delete(&User{}, "id = ?", userID).Error
}
