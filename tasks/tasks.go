package tasks

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	Pending    TaskStatus = "Pending"
	InProgress TaskStatus = "InProgress"
	Completed  TaskStatus = "Completed"
)

type Task struct {
	ID          uint       `gorm:"primaryKey;autoIncrement"`
	UserId      int64      `gorm:"index"`
	Description string     `gorm:"type:text;not null"`
	Status      TaskStatus `gorm:"type:varchar(20);default:'Pending'"`
	CreatedAt   time.Time
}

func AddTask(userId int64, description string, db *gorm.DB) (*Task, error) {
	task := &Task{
		UserId:      userId,
		Description: description,
		Status:      Pending,
		CreatedAt:   time.Now(),
	}

	if err := db.Create(&task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func UpdateTaskStatus(taskID uint, userId int64, updatedStatus TaskStatus, db *gorm.DB) error {
	if updatedStatus != Pending && updatedStatus != InProgress && updatedStatus != Completed {
		return errors.New("invalid status")
	}

	var task Task
	if err := db.First(&task, "id = ? AND user_id = ?", taskID, userId).Error; err != nil {
		return errors.New("task not found")
	}

	task.Status = updatedStatus
	if err := db.Save(&task).Error; err != nil {
		return err
	}

	return nil
}

func ListTasks(userId int64, db *gorm.DB) ([]Task, error) {
	var tasks []Task

	if err := db.Where("user_id = ?", userId).Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func DeleteTask(taskID uint, userId int64, db *gorm.DB) error {
	if err := db.Where("id = ? AND user_id = ?", taskID, userId).Delete(&Task{}).Error; err != nil {
		return err
	}

	return nil
}

func EditingTask(taskID uint, userId int64, newDescription string, db *gorm.DB) error {
	if newDescription == "" {
		return errors.New("description cannot be empty")
	}

	var task Task
	if err := db.First(&task, "id = ? AND user_id = ?", taskID, userId).Error; err != nil {
		return errors.New("task not found")
	}

	task.Description = newDescription
	if err := db.Save(&task).Error; err != nil {
		return errors.New("failed to update task description")
	}

	return nil
}
