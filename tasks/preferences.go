package tasks

import (
	"gorm.io/gorm"
)

type UserPreference struct {
	ID              uint  `gorm:"primaryKey;autoIncrement"`
	UserId          int64 `gorm:"index;unique;not null"`
	EnableReminders bool  `gorm:"default:true"`
}

func GetOrCreatePreference(userId int64, db *gorm.DB) (*UserPreference, error) {
	var preference UserPreference
	err := db.FirstOrCreate(&preference, UserPreference{UserId: userId}).Error
	return &preference, err
}

func UpdateUserPreference(userId int64, enableReminders bool, db *gorm.DB) error {
	var preference UserPreference
	err := db.First(&preference, "user_id = ?", userId).Error
	if err != nil {
		return err
	}

	preference.EnableReminders = enableReminders
	return db.Save(&preference).Error
}
