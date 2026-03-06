package handlers

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, error error, status int) {
	w.WriteHeader(status)
	bodyWithError, err := json.Marshal(error)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	} else {
		w.Write(bodyWithError)
	}
}
