package v1

import (
	"aggregator/internal/database"
	"aggregator/utils"
	"context"
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

		utils.RespondWithError(w, http.StatusInternalServerError, utils.InternalServerError)
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
	user, err := db.CreateUser(context.Background(), u)
	if err != nil {
		log.Print("Error creating user: ", err)

		utils.RespondWithError(w, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, user)
}
