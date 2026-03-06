package features

import (
	"booking-service/tools"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func InsertRow(ctx context.Context, conn *pgx.Conn, booking *tools.Booking) error {
	query := `INSERT INTO bookings (place_id, user_name, user_phone, start_time, end_time)
	 		VALUES ($1, $2, $3, $4, $5);`

	_, err := conn.Exec(
		ctx,
		query,
		booking.PlaceID,
		booking.UserName,
		booking.UserPhone,
		booking.StartTime,
		booking.EndTime,
	)
	fmt.Println(err)
	return err
}

func SelectAll(ctx context.Context, conn *pgx.Conn) ([]tools.Booking, error) {
	var booking tools.Booking
	var bookings []tools.Booking

	query := `SELECT place_id, user_name, user_phone, start_time, end_time FROM bookings`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		fmt.Println(err)
		return []tools.Booking{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&booking.PlaceID,
			&booking.UserName,
			&booking.UserPhone,
			&booking.StartTime,
			&booking.EndTime,
		)

		if err != nil {
			fmt.Println(err)
			return []tools.Booking{}, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}
