package v1

import (
	"aggregator/utils"
	"errors"
	"net/http"
)

var internalServerError = errors.New("Internal Server Error")

type statusResponse struct {
	Status string `json:"status,omitempty"`
}

func Readiness(w http.ResponseWriter, r *http.Request) {
	res := statusResponse{Status: "ok"}

	utils.RespondWithJson(w, http.StatusOK, res)
	return
}

func Err(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, http.StatusInternalServerError, internalServerError)
}
