package bots

import (
	"context"
	"errors"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Bot(ctx context.Context, conn *pgxpool.Pool) error {
	userID, err := strconv.Atoi(os.Getenv("ADMIN_ID"))
	if err != nil {
		return errors.New("invalid ADMIN_ID")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		return err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.From.ID != int64(userID) {
			continue
		}

		err := Handler(ctx, conn, bot, &update)
		if err != nil {
			return err
		}

	}
	return nil
}
