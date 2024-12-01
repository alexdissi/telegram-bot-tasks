package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"telegram-task-bot/tasks"
	"time"

	"github.com/go-telegram/bot"
	"gorm.io/gorm"
)

func StartNotificationWorker(ctx context.Context, db *gorm.DB, token string) {
	go func() {
		opts := []bot.Option{}
		b, err := bot.New(token, opts...)
		if err != nil {
			log.Fatalf("Error to create bot: %v", err)
		}

		for {
			select {
			case <-ctx.Done():
				log.Println("Notification stopped.")
				return
			default:
				err := sendPendingNotifications(ctx, db, b)
				if err != nil {
					log.Printf("Error sending notifications: %v", err)
				}
				time.Sleep(1 * time.Hour)
			}
		}
	}()
}

func TasksByUser(taskList []tasks.Task) map[int64][]string {
	userTasks := make(map[int64][]string)

	for _, task := range taskList {
		userTasks[task.UserId] = append(userTasks[task.UserId], task.Description)
	}

	return userTasks
}

func sendPendingNotifications(ctx context.Context, db *gorm.DB, b *bot.Bot) error {
	var taskList []tasks.Task

	err := db.Where("status = ? AND created_at <= ?", tasks.Pending, time.Now().Add(-24*time.Hour)).Find(&taskList).Error
	if err != nil {
		return fmt.Errorf("failed to fetch tasks: %w", err)
	}

	if len(taskList) == 0 {
		log.Println("No pending tasks found for notifications.")
		return nil
	}

	userTasks := TasksByUser(taskList)

	for userID, descriptions := range userTasks {
		var userPreference tasks.UserPreference
		err := db.First(&userPreference, "user_id = ?", userID).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				userPreference = tasks.UserPreference{
					UserId:          userID,
					EnableReminders: true,
				}
				_ = db.Create(&userPreference)
			} else {
				log.Printf("Error fetching user preferences for user %d: %v", userID, err)
				continue
			}
		}

		if !userPreference.EnableReminders {
			log.Printf("Reminders are disabled for user %d.", userID)
			continue
		}

		message := fmt.Sprintf("⏰ You have pending tasks:\n- %s", strings.Join(descriptions, "\n- "))
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   message,
		})
		if err != nil {
			log.Printf("Failed to send notification to user %d: %v", userID, err)
		} else {
			log.Printf("Notification sent to user %d", userID)
		}
	}

	return nil
}
