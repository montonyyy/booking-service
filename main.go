package main

import (
	bots "booking-service/features/bots/tgbot"
	"booking-service/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := &http.Server{Addr: ":8080"}

	connection, err := pgxpool.New(ctx, os.Getenv("CONN_STRING"))
	var conn = &handlers.Conn{
		Conn: connection,
		Ctx:  ctx,
	}

	if err != nil {
		log.Println("Ошибка при подключении к БД")
		panic(err)
	}
	defer connection.Close()

	http.HandleFunc("/booking", conn.SqlHandler)

	go func() {
		err := bots.Bot(ctx, connection)
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		err = http.ListenAndServe(":8080", nil)
		log.Println(err)
	}()

	<-ctx.Done()
	log.Println("Завершение работы")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Сервер остановлен")

}
