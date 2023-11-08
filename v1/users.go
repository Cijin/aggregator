package v1

import (
	"aggregator/internal/database"
	"aggregator/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
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
