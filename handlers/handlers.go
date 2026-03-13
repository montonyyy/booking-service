package handlers

import (
	features "booking-service/features/sql"
	"booking-service/tools"
	"context"
	"encoding/json"
	"io"
	"log/slog"
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
				slog.Error(err.Error())
			}
		}

	case http.MethodPost:
		var booking *tools.Booking

		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(body, &booking); err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		if err := features.InsertRow(c.Ctx, c.Conn, booking); err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(body); err != nil {
			slog.Error(err.Error())
		}

	case http.MethodDelete:
		var booking *tools.Booking

		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(body, &booking); err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		if err := features.DeleteRow(c.Ctx, c.Conn, booking); err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodPatch:
		var booking *tools.Booking

		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(body, &booking); err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		if err := features.UpdateRow(c.Ctx, c.Conn, booking); err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func writeError(w http.ResponseWriter, error error, status int) {
	w.WriteHeader(status)
	bodyWithError, err := json.Marshal(error)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
	if _, err := w.Write(bodyWithError); err != nil {
		slog.Error(err.Error())
	}
}
