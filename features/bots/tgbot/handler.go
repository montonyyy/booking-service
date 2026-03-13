package bots

import (
	features "booking-service/features/sql"
	"booking-service/tools"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Handler(ctx context.Context, conn *pgxpool.Pool, bot *tgbotapi.BotAPI, u *tgbotapi.Update, updates *tgbotapi.UpdatesChannel) error {
	switch u.Message.Text {
	case "/start":
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, `
			/list - показать бронирования
			/add - добавить бронирование
			/del - удалить бронирование
			/rep - изменить бронирование
			`)

		if _, err := bot.Send(msg); err != nil {
			return err
		}

	case "/list":
		var lines []string
		table, err := features.SelectAll(ctx, conn)
		if err != nil {
			return err
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
			return err
		}

	case "/add":
		booking := &tools.Booking{}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, `
					Чтобы добавить запись, введите следующие значения:
					{place_id} {user_name} {user_phone} {start_time} {end_time}
					`)
		if _, err := bot.Send(msg); err != nil {
			return err
		}

		nextUpdate := WaitNextUpdate(*updates)
		if nextUpdate == nil {
			return errors.New("time for data entry has expired")
		}

		if len(strings.Split(nextUpdate.Message.Text, " ")) == 5 {
			if _, err := fmt.Sscanf(nextUpdate.Message.Text, "%d %s %s %s %s",
				&booking.PlaceID,
				&booking.UserName,
				&booking.UserPhone,
				&booking.StartTime,
				&booking.EndTime,
			); err != nil {
				return err
			}
		} else {
			return errors.New("incorrect data entry format")
		}

		if err := features.InsertRow(ctx, conn, booking); err != nil {
			return err
		}

	case "/del":
		booking := &tools.Booking{}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, `
				Чтобы удалить запись, введите следующее значение:
				{id}
				`)

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		nextUpdate := WaitNextUpdate(*updates)
		if nextUpdate == nil {
			return errors.New("time for data entry has expired")
		}

		if len(strings.Split(nextUpdate.Message.Text, " ")) == 1 {
			if _, err := fmt.Sscanf(nextUpdate.Message.Text, "%d",
				&booking.ID,
			); err != nil {
				return err
			}
		} else {
			return errors.New("incorrect data entry format")
		}

		if err := features.DeleteRow(ctx, conn, booking); err != nil {
			return err
		}

	case "/rep":
		booking := &tools.Booking{}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, `
				Чтобы изменить запись, введите следующее значение:
				{ID изм. записи} {Новая нач. дата} {Новая кон. дата}
				`)
		if _, err := bot.Send(msg); err != nil {
			return err
		}

		nextUpdate := WaitNextUpdate(*updates)
		if nextUpdate == nil {
			return errors.New("time for data entry has expired")
		}

		if len(strings.Split(nextUpdate.Message.Text, " ")) == 3 {
			if _, err := fmt.Sscanf(nextUpdate.Message.Text, "%d %s %s",
				&booking.ID, &booking.StartTime, &booking.EndTime,
			); err != nil {
				return err
			}
		} else {
			return errors.New("incorrect data entry format")
		}

		if err := features.UpdateRow(ctx, conn, booking); err != nil {
			return err
		}
	}
	return nil
}

func WaitNextUpdate(updates tgbotapi.UpdatesChannel) *tgbotapi.Update {
	select {
	case update := <-updates:
		return &update
	case <-time.After(60 * time.Second):
		return nil
	}
}
