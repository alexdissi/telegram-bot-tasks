package tasks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDb(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&Task{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestAddTask(t *testing.T) {
	db := setupTestDb(t)
	task, err := AddTask(12345, "Test Task", db)
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	assert.Equal(t, uint(1), task.ID)
	assert.Equal(t, int64(12345), task.UserId)
	assert.Equal(t, "Test Task", task.Description)
	assert.Equal(t, Pending, task.Status)
	assert.WithinDuration(t, time.Now(), task.CreatedAt, time.Second)
}

func TestListTasks(t *testing.T) {
	db := setupTestDb(t)
	_, _ = AddTask(12345, "Task 1", db)
	_, _ = AddTask(12345, "Task 2", db)
	_, _ = AddTask(67890, "Other User Task", db)

	tasks, err := ListTasks(12345, db)
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	assert.Len(t, tasks, 2)
	assert.Equal(t, "Task 1", tasks[0].Description)
	assert.Equal(t, "Task 2", tasks[1].Description)
}

func TestUpdateTaskStatus(t *testing.T) {
	db := setupTestDb(t)
	task, _ := AddTask(12345, "Task to Update", db)

	err := UpdateTaskStatus(task.ID, task.UserId, Completed, db)
	if err != nil {
		t.Fatalf("Failed to update task status: %v", err)
	}

	var updatedTask Task
	err = db.First(&updatedTask, task.ID).Error
	if err != nil {
		t.Fatalf("Failed to fetch updated task: %v", err)
	}

	assert.Equal(t, Completed, updatedTask.Status)
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDb(t)
	task, _ := AddTask(12345, "Task to Delete", db)

	err := DeleteTask(task.ID, task.UserId, db)
	assert.NoError(t, err)

	var deletedTask Task
	err = db.First(&deletedTask, task.ID).Error
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestEditingTask(t *testing.T) {
	db := setupTestDb(t)
	task, _ := AddTask(12345, "Task to Edit", db)

	err := EditingTask(task.ID, task.UserId, "Edited Task", db)
	assert.NoError(t, err)

	var updatedTask Task
	err = db.First(&updatedTask, task.ID).Error
	if err != nil {
		t.Fatalf("Failed to fetch updated task: %v", err)
	}

	assert.Equal(t, "Edited Task", updatedTask.Description)
}
