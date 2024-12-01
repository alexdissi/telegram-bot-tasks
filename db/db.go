package db

import (
	"log"
	"telegram-task-bot/tasks"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("telegram_tasks.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error with database connection : %v", err)
	}

	err = db.AutoMigrate(&tasks.Task{}, &tasks.UserPreference{})
	if err != nil {
		log.Fatalf("Erreur to migrate model : %v", err)
	}

	log.Println("Connexion bdd OK !")
	return db
}
