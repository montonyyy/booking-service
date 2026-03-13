package main

import (
	bots "booking-service/features/bots/tgbot"
	"booking-service/handlers"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := os.Getenv("SERVER_PORT")
	if string(port[0]) != ":" {
		port = ":" + port
	}

	server := &http.Server{Addr: port}
	connection, err := pgxpool.New(ctx, os.Getenv("CONN_STRING"))

	if err != nil {
		slog.Error(err.Error())
	}
	defer connection.Close()

	listener, err := pgx.Connect(ctx, os.Getenv("CONN_STRING"))
	if err != nil {
		slog.Error(err.Error())
	}
	defer listener.Close(ctx)

	if _, err := listener.Exec(ctx, "LISTEN updates"); err != nil {
		slog.Error(err.Error())
	}

	var conn = &handlers.Conn{
		Conn: connection,
		Ctx:  ctx,
	}

	http.HandleFunc("/booking", conn.SqlHandler)

	go func() {
		err := bots.Bot(ctx, connection, listener)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	go func() {
		err = http.ListenAndServe(port, nil)

		slog.Info(err.Error())
	}()

	<-ctx.Done()
	slog.Info("Завершение работы")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error(err.Error())
	}

	slog.Info("Сервер остановлен")
}
