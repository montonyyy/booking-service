package bots

import (
	features "booking-service/features/sql"
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Handler(ctx context.Context, conn *pgxpool.Pool, bot *tgbotapi.BotAPI, u *tgbotapi.Update) error {

	var err error = nil

	switch u.Message.Text {

	case "/start":
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, `
			/list - показать бронирования
			`)
		bot.Send(msg)

	case "/list":
		table, err := features.SelectAll(ctx, conn)
		if err != nil {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка при получении данных")
			bot.Send(msg)
			break
		}
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, fmt.Sprintf("%v", table))
		bot.Send(msg)

		/*
			case "/add":
				var booking *tools.Booking

				booking, err := // --> TODO
				if err != nil {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка при чтении сообщения")
					bot.Send(msg)
					break
				}

				err = features.InsertRow(ctx, conn, booking)
				if err != nil {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка при вставки данных")
					bot.Send(msg)
					break
				}
				msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Данные записаны")
				bot.Send(msg)

			case "/del":
				var booking *tools.Booking

				booking, err := // --> TODO
				if err != nil {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка при чтении сообщения")
					bot.Send(msg)
					break
				}

				err = features.DeleteRow(ctx, conn, booking)
				if err != nil {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка при удалении данных")
					bot.Send(msg)
					break
				}
				msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Данные удалены")
				bot.Send(msg)
		*/
	}

	// listen --> TODO

	return err
}
