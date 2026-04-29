package models

import (
	"task_manager/db"
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

func GetTasksByProjectId(projectID uuid.UUID) ([]Task, error) {
	var tasks []Task
	if err := db.DB.Where("project_id = ?", projectID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// Fetch all tasks assigned to a user
func GetTasksByUserId(userID uuid.UUID) ([]Task, error) {
	var tasks []Task
	if err := db.DB.Where("assigned_to = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// Create a new task
func CreateTask(task *Task) error {
	return db.DB.Create(task).Error
}

// Get a task by ID
func GetTaskById(taskID uuid.UUID) (*Task, error) {
	var task Task
	if err := db.DB.First(&task, "id = ?", taskID).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// Update task status or due date
func UpdateTask(taskID uuid.UUID, updates map[string]interface{}) error {
	return db.DB.Model(&Task{}).Where("id = ?", taskID).Updates(updates).Error
}

// Delete a task
func DeleteTask(taskID uuid.UUID) error {
	return db.DB.Delete(&Task{}, "id = ?", taskID).Error
}

// Get overdue tasks (example query)
func GetOverdueTasks(now time.Time) ([]Task, error) {
	var tasks []Task
	if err := db.DB.Where("due_date < ? AND status != ?", now, "done").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
