package handlers

import (
	features "booking-service/features/sql"
	"booking-service/tools"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Conn struct {
	Conn *pgxpool.Pool
	Ctx  context.Context
}

func (c *Conn) SqlHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		table, err := features.SelectAll(c.Ctx, c.Conn)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		tableMarshal, err := json.Marshal(table)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(tableMarshal)
		}

	case http.MethodPost:
		var booking *tools.Booking

		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &booking)
		if err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		err = features.InsertRow(c.Ctx, c.Conn, booking)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)

	case http.MethodDelete:
		var booking *tools.Booking

		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &booking)
		if err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		err = features.DeleteRow(c.Ctx, c.Conn, booking)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)

	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"unsupported method"}`))
	}
}

func writeError(w http.ResponseWriter, error error, status int) {
	w.WriteHeader(status)
	bodyWithError, err := json.Marshal(error)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Write(bodyWithError)

}
