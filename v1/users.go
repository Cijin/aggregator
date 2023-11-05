package v1

import (
	"aggregator/internal/database"
	"aggregator/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (v *v1) CreateUser(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Name string `json:"name,omitempty"`
	}

	b := &reqBody{}
	d := json.NewDecoder(r.Body)

	err := d.Decode(b)
	if err != nil {
		log.Print("Json decode error:", err)

		utils.RespondWithInternalServerError(w)
		return
	}

	if b.Name == "" {
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("name cannot be empty"))
		return
	}

	db := v.Db

	u := database.CreateUserParams{
		Name:      b.Name,
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user, err := db.CreateUser(r.Context(), u)
	if err != nil {
		log.Print("Error creating user: ", err)

		utils.RespondWithInternalServerError(w)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, user)
}

func (v *v1) GetUserByApiKey(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	headerComponents := strings.Split(authHeader, " ")

	if len(headerComponents) != 2 || headerComponents[0] != "ApiKey" {
		utils.RespondWithError(w, http.StatusUnauthorized, errors.New("invalid auth header"))
		return
	}

	if headerComponents[1] == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, errors.New("missing api key"))
		return
	}

	apiKey := headerComponents[1]

	user, err := v.Db.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusUnauthorized, errors.New("invalid api key"))
			return
		}

		log.Print("Error getting user by api key: ", err)

		utils.RespondWithInternalServerError(w)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, user)
}
