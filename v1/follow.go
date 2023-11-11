package v1

import (
	"aggregator/internal/database"
	"aggregator/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (v *v1) Follow(w http.ResponseWriter, r *http.Request, u database.User) {
	type reqBody struct {
		FeedId string `json:"feed_id,omitempty"`
	}

	req := reqBody{}
	d := json.NewDecoder(r.Body)

	err := d.Decode(&req)
	if err != nil {
		log.Print("err decoding json:", err)

		utils.RespondWithInternalServerError(w)
		return
	}

	if req.FeedId == "" {
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("feed id cannot be empty"))
		return
	}

	feedId, err := uuid.Parse(req.FeedId)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("feed id is not a valid"))
		return
	}

	followParam := database.FollowParams{
		ID:        uuid.New(),
		FeedID:    feedId,
		UserID:    u.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	f, err := v.Db.Follow(r.Context(), followParam)
	if err != nil {
		log.Print("error creating follow: ", err)

		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, f)
}

func (v *v1) Unfollow(w http.ResponseWriter, r *http.Request, u database.User) {
	feedFollowId := chi.URLParam(r, "feedFollowId")

	feedFollowUuid, err := uuid.Parse(feedFollowId)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("feed id is not a valid"))
		return
	}

	err = v.Db.DeleteFeedFollow(r.Context(), feedFollowUuid)
	if err != nil {
		log.Print("error creating follow: ", err)

		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, "ok")
}

func (v *v1) GetFeed(w http.ResponseWriter, r *http.Request, u database.User) {
	feeds, err := v.Db.GetUserFeedFollow(r.Context(), u.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithJson(w, http.StatusOK, "user is not following any feed")
		}
		log.Print("error creating follow: ", err)

		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, feeds)
}
