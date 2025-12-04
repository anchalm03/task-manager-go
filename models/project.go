package models

import (
	"task_manager/db"

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

func (Project) TableName() string {
	return "projects"
}

func GetAllProjectsByUserId(userID uuid.UUID) ([]Project, error) {
	var projects []Project
	if err := db.DB.Where("owner_id = ?", userID).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func CreateProject(project *Project) error {
	return db.DB.Create(project).Error
}

func GetProjectById(projectID uuid.UUID) (*Project, error) {
	var project Project
	if err := db.DB.First(&project, "id = ?", projectID).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func DeleteProject(projectID uuid.UUID) error {
	return db.DB.Delete(&Project{}, "id = ?", projectID).Error
}
