package v1

import (
	"aggregator/internal/database"
	"aggregator/utils"
	"database/sql"
	"errors"
	"log"
	"net/http"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (v *v1) GetUserByApiKey(r *http.Request) (database.User, error) {
	apiKey, err := utils.GetApiKey(r)
	if err != nil {
		return database.User{}, err
	}

	user, err := v.Db.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return database.User{}, errors.New("invalid api key")
		}

		log.Print("Error getting user by api key: ", err)
		return database.User{}, errors.New("internal server error")
	}

	return user, nil
}

func (v *v1) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := v.GetUserByApiKey(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, err)
			return
		}

		handler(w, r, user)
	}
}
