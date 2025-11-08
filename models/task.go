package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Id          uuid.UUID  `gorm:"column:id;default:gen_random_uuid();primaryKey"`
	Title       string     `gorm:"column:title;not null"`
	Description string     `gorm:"column:description"`
	Status      string     `gorm:"column:status;default:'todo'"`
	DueDate     *time.Time `gorm:"column:due_date"`
	AssignedTo  *uuid.UUID `gorm:"column:assigned_to"` // nullable
	ProjectID   uuid.UUID  `gorm:"column:project_id"`  // required, CASCADE delete on project
}
