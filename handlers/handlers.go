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
			if _, err := w.Write(tableMarshal); err != nil {
				log.Panic(err)
			}
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
		if _, err := w.Write(body); err != nil {
			log.Panic(err)
		}

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
		if _, err := w.Write(body); err != nil {
			log.Panic(err)
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(`{"error":"unsupported method"}`)); err != nil {
			log.Panic(err)
		}
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
	if _, err := w.Write(bodyWithError); err != nil {
		log.Panic(err)
	}
}
