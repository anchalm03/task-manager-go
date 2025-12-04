package models

import (
	"task_manager/db"

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

func GetCommentsByTaskId(taskID uuid.UUID) ([]Comment, error) {
	var comments []Comment
	if err := db.DB.Where("task_id = ?", taskID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// Add comment to a task
func CreateComment(comment *Comment) error {
	return db.DB.Create(comment).Error
}

// Delete comment
func DeleteComment(commentID uuid.UUID) error {
	return db.DB.Delete(&Comment{}, "id = ?", commentID).Error
}
