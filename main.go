package main

import (
	"booking-service/handlers"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	connection, err := pgx.Connect(ctx, os.Getenv("CONN_STRING"))
	var conn = &handlers.Conn{
		Conn: connection,
		Ctx:  ctx,
	}

	if err != nil {
		fmt.Println("Ошибка при подключении к БД")
		panic(err)
	}
	defer connection.Close(ctx)

	http.HandleFunc("/booking", conn.SqlHandler)

	_ = http.ListenAndServe(":8080", nil)

}
