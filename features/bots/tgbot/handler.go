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
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

	case "/list":
		var lines []string
		table, err := features.SelectAll(ctx, conn)
		if err != nil {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Ошибка при получении данных")
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			break
		}

		for _, line := range table {
			str := fmt.Sprintf("ID: %d, Место: %d, Имя: %s, Телефон: %s, Нач. дата: %s, Кон. дата: %s \n",
				line.ID,
				line.PlaceID,
				line.UserName,
				line.UserPhone,
				line.StartTime,
				line.EndTime,
			)
			lines = append(lines, str)
		}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, strings.Join(lines, "\n"))
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

	case "/add":
		booking := &tools.Booking{}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, `
					Чтобы добавить запись, введите следующие значения:
					{place_id} {user_name} {user_phone} {start_time} {end_time}
					`)
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

		nextUpdate := WaitNextUpdate(*updates)
		if nextUpdate == nil {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Время ввода данных истекло. Повторите попытку")
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			break
		}

		if len(strings.Split(nextUpdate.Message.Text, " ")) == 5 {

			if _, err := fmt.Sscanf(nextUpdate.Message.Text, "%d %s %s %s %s",
				&booking.PlaceID,
				&booking.UserName,
				&booking.UserPhone,
				&booking.StartTime,
				&booking.EndTime,
			); err != nil {
				log.Panic(err)
			}
		} else {
			log.Println("Ошибка при чтении данных")
			msg := tgbotapi.NewMessage(nextUpdate.Message.Chat.ID, "Ошибка при чтении данных")
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			break
		}

		err = features.InsertRow(ctx, conn, booking)
		if err != nil {
			msg := tgbotapi.NewMessage(nextUpdate.Message.Chat.ID, "Ошибка при вставки данных")
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			break
		}

	case "/del":
		booking := &tools.Booking{}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, `
				Чтобы удалить запись, введите следующее значение:
				{id}
				`)
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

		nextUpdate := WaitNextUpdate(*updates)
		if nextUpdate == nil {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Время ввода данных истекло. Повторите попытку")
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			break
		}

		if len(strings.Split(nextUpdate.Message.Text, " ")) == 1 {

			if _, err := fmt.Sscanf(nextUpdate.Message.Text, "%d",
				&booking.ID,
			); err != nil {
				log.Panic(err)
			}
		} else {
			log.Println("Ошибка при чтении данных")
			msg := tgbotapi.NewMessage(nextUpdate.Message.Chat.ID, "Ошибка при чтении данных")
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			break
		}

		err = features.DeleteRow(ctx, conn, booking)
		if err != nil {
			msg := tgbotapi.NewMessage(nextUpdate.Message.Chat.ID, "Ошибка при удалении данных")
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			break
		}
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
