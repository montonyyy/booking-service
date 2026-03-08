package features

import (
	"booking-service/tools"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertRow(ctx context.Context, conn *pgxpool.Pool, booking *tools.Booking) error {
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
	if err != nil {
		log.Println(err)
	} else {
		conn.Exec(ctx, "NOTIFY updates, $1", fmt.Sprintf("%s, %s", booking.UserName, booking.UserPhone))
	}

	return err
}

func SelectAll(ctx context.Context, conn *pgxpool.Pool) ([]tools.Booking, error) {
	var booking tools.Booking
	var bookings []tools.Booking

	query := `SELECT id, place_id, user_name, user_phone, start_time, end_time FROM bookings`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		log.Println(err)
		return []tools.Booking{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&booking.ID,
			&booking.PlaceID,
			&booking.UserName,
			&booking.UserPhone,
			&booking.StartTime,
			&booking.EndTime,
		)

		if err != nil {
			log.Println(err)
			return []tools.Booking{}, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func DeleteRow(ctx context.Context, conn *pgxpool.Pool, booking *tools.Booking) error {
	query := `DELETE FROM bookings WHERE id = $1;`

	_, err := conn.Exec(ctx, query, booking.ID)
	if err != nil {
		log.Println(err)
	}
	return err

}
