package bots

import (
	features "booking-service/features/sql"
	"booking-service/tools"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Handler(ctx context.Context, conn *pgxpool.Pool, bot *tgbotapi.BotAPI, u *tgbotapi.Update, updates *tgbotapi.UpdatesChannel) error {

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

	case "/add":
		booking := &tools.Booking{}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, `
					Чтобы добавить запись, введите следующие значения:
					{place_id} {user_name} {user_phone} {start_time} {end_time}
					`)
		bot.Send(msg)

		nextUpdate := WaitNextUpdate(*updates)

		if len(strings.Split(nextUpdate.Message.Text, " ")) == 5 {

			fmt.Sscanf(nextUpdate.Message.Text, "%d %s %s %s %s",
				&booking.PlaceID,
				&booking.UserName,
				&booking.UserPhone,
				&booking.StartTime,
				&booking.EndTime,
			)
		} else {
			log.Println("Ошибка при чтении данных")
			msg := tgbotapi.NewMessage(nextUpdate.Message.Chat.ID, "Ошибка при чтении данных")
			bot.Send(msg)
			break
		}

		err = features.InsertRow(ctx, conn, booking)
		if err != nil {
			msg := tgbotapi.NewMessage(nextUpdate.Message.Chat.ID, "Ошибка при вставки данных")
			bot.Send(msg)
			break
		}
		msg = tgbotapi.NewMessage(nextUpdate.Message.Chat.ID, "Данные записаны")
		bot.Send(msg)

	case "/del":
		booking := &tools.Booking{}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, `
				Чтобы удалить запись, введите следующее значение:
				{id}
				`)
		bot.Send(msg)

		nextUpdate := WaitNextUpdate(*updates)

		if len(strings.Split(nextUpdate.Message.Text, " ")) == 1 {

			fmt.Sscanf(nextUpdate.Message.Text, "%d",
				&booking.ID,
			)
		} else {
			log.Println("Ошибка при чтении данных")
			msg := tgbotapi.NewMessage(nextUpdate.Message.Chat.ID, "Ошибка при чтении данных")
			bot.Send(msg)
			break
		}

		err = features.DeleteRow(ctx, conn, booking)
		if err != nil {
			msg := tgbotapi.NewMessage(nextUpdate.Message.Chat.ID, "Ошибка при удалении данных")
			bot.Send(msg)
			break
		}
		msg = tgbotapi.NewMessage(nextUpdate.Message.Chat.ID, "Данные удалены")
		bot.Send(msg)
	}

	return err
}

func WaitNextUpdate(updates tgbotapi.UpdatesChannel) *tgbotapi.Update {
	select {
	case update := <-updates:
		return &update
	case <-time.After(60 * time.Second):
		return nil
	}
}
