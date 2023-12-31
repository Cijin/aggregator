package v1

import (
	"aggregator/internal/database"
	"aggregator/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (v *v1) CreateFeed(w http.ResponseWriter, r *http.Request, u database.User) {
	type reqBody struct {
		Name string `json:"name,omitempty"`
		Url  string `json:"url,omitempty"`
	}

	d := json.NewDecoder(r.Body)
	req := &reqBody{}

	err := d.Decode(req)
	if err != nil {
		utils.RespondWithInternalServerError(w)
		return
	}

	f := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      req.Name,
		Url:       req.Url,
		UserID:    u.ID,
	}

	feed, err := v.Db.CreateFeed(r.Context(), f)
	if err != nil {
		utils.RespondWithInternalServerError(w)
		return
	}

	// automatically create follow
	_, err = v.Db.Follow(r.Context(), database.FollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    u.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Print("error creating follow:", err)
	}

	utils.RespondWithJson(w, http.StatusOK, feed)
}

func (v *v1) ListFeed(w http.ResponseWriter, r *http.Request) {
	feeds, err := v.Db.ListFeed(r.Context())
	if err != nil {
		utils.RespondWithInternalServerError(w)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, feeds)
}
