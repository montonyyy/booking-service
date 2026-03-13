package features

import (
	"booking-service/tools"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertRow(ctx context.Context, conn *pgxpool.Pool, booking *tools.Booking) error {
	query := `INSERT INTO bookings (place_id, user_name, user_phone, start_time, end_time)
	 		VALUES ($1, $2, $3, $4, $5);`

	if _, err := conn.Exec(
		ctx,
		query,
		booking.PlaceID,
		booking.UserName,
		booking.UserPhone,
		booking.StartTime,
		booking.EndTime,
	); err != nil {
		return err
	}

	payload := fmt.Sprintf("Данные добавлены: %s, %s", booking.UserName, booking.UserPhone)
	payload = strings.ReplaceAll(payload, "'", "''")

	query = fmt.Sprintf("NOTIFY updates, '%s'", payload)

	if _, err := conn.Exec(ctx, query); err != nil {
		return err
	}
	return nil
}

func SelectAll(ctx context.Context, conn *pgxpool.Pool) ([]tools.Booking, error) {
	var booking tools.Booking
	var bookings []tools.Booking

	query := `SELECT id, place_id, user_name, user_phone, start_time, end_time FROM bookings`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return []tools.Booking{}, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&booking.ID,
			&booking.PlaceID,
			&booking.UserName,
			&booking.UserPhone,
			&booking.StartTime,
			&booking.EndTime,
		); err != nil {
			return []tools.Booking{}, err
		}

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func DeleteRow(ctx context.Context, conn *pgxpool.Pool, booking *tools.Booking) error {
	query := `DELETE FROM bookings WHERE id = $1;`

	if _, err := conn.Exec(ctx, query, booking.ID); err != nil {
		return err
	}
	if _, err := conn.Exec(ctx, "NOTIFY updates, 'Данные удалены'"); err != nil {
		return err
	}
	return nil
}

func UpdateRow(ctx context.Context, conn *pgxpool.Pool, booking *tools.Booking) error {
	query := `UPDATE bookings
				SET start_time=$1, end_time=$2
				WHERE id=$3;`

	if _, err := conn.Exec(ctx, query, booking.StartTime, booking.EndTime, booking.ID); err != nil {
		return err
	}
	if _, err := conn.Exec(ctx, "NOTIFY updates, 'Данные обновлены'"); err != nil {
		return err
	}
	return nil
}
