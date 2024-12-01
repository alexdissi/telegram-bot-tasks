package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"telegram-task-bot/bot"
	"telegram-task-bot/db"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error: .env file not found")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("Error: Telegram bot token is missing")
	}

	database := db.ConnectDatabase()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Println("Starting Telegram bot...")
	err = bot.Start(ctx, token, database)
	if err != nil {
		log.Fatalf("Error starting the bot: %v", err)
	}

	log.Println("Starting notification service...")
	bot.StartNotificationWorker(ctx, database, token)

	log.Println("Bot is running. Press Ctrl+C to stop.")

	<-ctx.Done()

	log.Println("Stopping the bot.")
}
