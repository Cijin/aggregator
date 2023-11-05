package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var InternalServerError = errors.New("internal server error")

type errRes struct {
	Error string `json:"error,omitempty"`
}

func RespondWithError(w http.ResponseWriter, status int, err error) {
	d := errRes{Error: err.Error()}

	RespondWithJson(w, status, d)
}

func RespondWithJson(w http.ResponseWriter, status int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Marshal error %s", err.Error())
		return
	}

	w.WriteHeader(status)
	w.Write(res)
}
