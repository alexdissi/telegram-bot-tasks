package utils

import (
	"context"
	"log"
	"time"

	"github.com/go-telegram/bot"
)

func RetryTelegramMessage(ctx context.Context, b *bot.Bot, params *bot.SendMessageParams, maxRetries int) error {
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		_, err = b.SendMessage(ctx, params)
		if err == nil {
			return nil
		}

		log.Printf("Error sending message (attempt %d/%d): %v", attempt, maxRetries, err)

		time.Sleep(2 * time.Second)
	}

	return err
}
