package v1

import (
	"aggregator/internal/database"
	"aggregator/utils"
	"errors"
	"net/http"
)

var internalServerError = errors.New("Internal Server Error")

type v1 struct {
	Db *database.Queries
}

func NewHandler(db *database.Queries) *v1 {
	return &v1{Db: db}
}

type statusResponse struct {
	Status string `json:"status,omitempty"`
}

func (v *v1) Readiness(w http.ResponseWriter, r *http.Request) {
	res := statusResponse{Status: "ok"}

	utils.RespondWithJson(w, http.StatusOK, res)
	return
}

func (v *v1) Err(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, http.StatusInternalServerError, internalServerError)
}
