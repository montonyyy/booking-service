package handlers

import (
	features "booking-service/features/sql"
	"booking-service/tools"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type Conn struct {
	Conn *pgx.Conn
	Ctx  context.Context
}

func (c *Conn) SqlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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
	} else if r.Method == http.MethodPost {
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
	} else if r.Method == http.MethodDelete {
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
	} else {
		writeError(w, errors.New("unsupported method"), http.StatusBadRequest)
		return
	}

}
