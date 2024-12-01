package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	utils "telegram-task-bot/error"
	"telegram-task-bot/tasks"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/gorm"
)

func Start(ctx context.Context, token string, db *gorm.DB) error {
	opts := []bot.Option{
		bot.WithDefaultHandler(DefaultHandler),
	}
	b, err := bot.New(token, opts...)
	if err != nil {
		return err
	}

	registerHandlers(b, db)
	go b.Start(ctx)

	return nil
}

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Command not recognized. Type /start to see available commands.",
	})
}

func registerHandlers(b *bot.Bot, db *gorm.DB) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		sendMessage(ctx, b, update.Message.Chat.ID, "Welcome to the task management bot! 🚀\n\nCommands:\n/add <task>\n/list\n/done <id>\n/delete <id>\n/edit <id> <new_description>\n/enable_reminders\n/disable_reminders")
	})

	b.RegisterHandler(bot.HandlerTypeMessageText, "/add", bot.MatchTypePrefix, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		handleAddTask(ctx, b, update, db)
	})

	b.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		handleListTasks(ctx, b, update, db)
	})

	b.RegisterHandler(bot.HandlerTypeMessageText, "/done", bot.MatchTypePrefix, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		handleUpdateTaskStatus(ctx, b, update, db, tasks.Completed)
	})

	b.RegisterHandler(bot.HandlerTypeMessageText, "/delete", bot.MatchTypePrefix, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		handleDeleteTask(ctx, b, update, db)
	})

	b.RegisterHandler(bot.HandlerTypeMessageText, "/edit", bot.MatchTypePrefix, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		handleEditTask(ctx, b, update, db)
	})

	b.RegisterHandler(bot.HandlerTypeMessageText, "/enable_reminders", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		handleEnableReminders(ctx, b, update, db)
	})

	b.RegisterHandler(bot.HandlerTypeMessageText, "/disable_reminders", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		handleDisableReminders(ctx, b, update, db)
	})
}

func handleAddTask(ctx context.Context, b *bot.Bot, update *models.Update, db *gorm.DB) {
	description := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/add"))
	if description == "" {
		sendMessage(ctx, b, update.Message.Chat.ID, "Provide a task description. Example: /add Buy groceries")
		return
	}

	_, err := tasks.GetOrCreatePreference(update.Message.Chat.ID, db)
	if err != nil {
		sendMessage(ctx, b, update.Message.Chat.ID, "Error adding task.")
		return
	}

	task, err := tasks.AddTask(update.Message.Chat.ID, description, db)
	if err != nil {
		sendMessage(ctx, b, update.Message.Chat.ID, "Error adding task.")
		return
	}

	sendMessage(ctx, b, update.Message.Chat.ID, fmt.Sprintf("Task [%v] added: %v", task.ID, task.Description))
}

func handleListTasks(ctx context.Context, b *bot.Bot, update *models.Update, db *gorm.DB) {
	taskList, err := tasks.ListTasks(update.Message.Chat.ID, db)
	if err != nil {
		sendMessage(ctx, b, update.Message.Chat.ID, "Error retrieving tasks.")
		return
	}

	if len(taskList) == 0 {
		sendMessage(ctx, b, update.Message.Chat.ID, "No tasks found.")
		return
	}

	var message strings.Builder
	message.WriteString("Your tasks:\n")
	for _, task := range taskList {
		message.WriteString(fmt.Sprintf("[%v] %v (%v)\n", task.ID, task.Description, task.Status))
	}

	sendMessage(ctx, b, update.Message.Chat.ID, message.String())
}

func handleUpdateTaskStatus(ctx context.Context, b *bot.Bot, update *models.Update, db *gorm.DB, newStatus tasks.TaskStatus) {
	idText := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/done"))
	taskId, err := strconv.ParseUint(idText, 10, 32)
	if err != nil {
		sendMessage(ctx, b, update.Message.Chat.ID, "Invalid task ID.")
		return
	}

	err = tasks.UpdateTaskStatus(uint(taskId), update.Message.Chat.ID, newStatus, db)
	if err != nil {
		sendMessage(ctx, b, update.Message.Chat.ID, fmt.Sprintf("Error: %v", err))
		return
	}

	sendMessage(ctx, b, update.Message.Chat.ID, fmt.Sprintf("Task [%v] marked as %v.", taskId, newStatus))
}

func handleDeleteTask(ctx context.Context, b *bot.Bot, update *models.Update, db *gorm.DB) {
	idText := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/delete"))
	taskId, err := strconv.ParseUint(idText, 10, 32)
	if err != nil {
		sendMessage(ctx, b, update.Message.Chat.ID, "Invalid task ID.")
		return
	}

	err = tasks.DeleteTask(uint(taskId), update.Message.Chat.ID, db)
	if err != nil {
		sendMessage(ctx, b, update.Message.Chat.ID, fmt.Sprintf("Error: %v", err))
		return
	}

	sendMessage(ctx, b, update.Message.Chat.ID, fmt.Sprintf("Task [%v] deleted.", taskId))
}

func handleEditTask(ctx context.Context, b *bot.Bot, update *models.Update, db *gorm.DB) {
	commandArgs := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/edit"))
	parts := strings.SplitN(commandArgs, " ", 2)
	if len(parts) < 2 {
		sendMessage(ctx, b, update.Message.Chat.ID, "Use /edit <task_id> <new_description>")
		return
	}

	taskID, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		sendMessage(ctx, b, update.Message.Chat.ID, "Invalid task ID.")
		return
	}

	newDescription := strings.TrimSpace(parts[1])
	if newDescription == "" {
		sendMessage(ctx, b, update.Message.Chat.ID, "Description cannot be empty.")
		return
	}

	err = tasks.EditingTask(uint(taskID), update.Message.Chat.ID, newDescription, db)
	if err != nil {
		sendMessage(ctx, b, update.Message.Chat.ID, fmt.Sprintf("Error: %v", err))
		return
	}

	sendMessage(ctx, b, update.Message.Chat.ID, fmt.Sprintf("Task [%v] updated to: %s", taskID, newDescription))
}

func handleEnableReminders(ctx context.Context, b *bot.Bot, update *models.Update, db *gorm.DB) {
	err := tasks.UpdateUserPreference(update.Message.Chat.ID, true, db)
	if err != nil {
		sendMessage(ctx, b, update.Message.Chat.ID, "Error enabling reminders.")
		return
	}
	sendMessage(ctx, b, update.Message.Chat.ID, "✅ Reminders enabled.")
}

func handleDisableReminders(ctx context.Context, b *bot.Bot, update *models.Update, db *gorm.DB) {
	err := tasks.UpdateUserPreference(update.Message.Chat.ID, false, db)
	if err != nil {
		sendMessage(ctx, b, update.Message.Chat.ID, "Error disabling reminders.")
		return
	}
	sendMessage(ctx, b, update.Message.Chat.ID, "❌ Reminders disabled.")
}

func sendMessage(ctx context.Context, b *bot.Bot, chatID int64, message string) {
	params := &bot.SendMessageParams{
		ChatID: chatID,
		Text:   message,
	}

	err := utils.RetryTelegramMessage(ctx, b, params, 3)
	if err != nil {
		log.Printf("Failed to send message to chat %d: %v", chatID, err)
	}
}
