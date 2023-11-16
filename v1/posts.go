package v1

import (
	"aggregator/internal/database"
	"aggregator/utils"
	"log"
	"net/http"
)

func (v *v1) ListPosts(w http.ResponseWriter, r *http.Request, u database.User) {
	param := database.ListsPostsParams{
		UserID: u.ID,
		Limit:  10,
	}

	posts, err := v.Db.ListsPosts(r.Context(), param)
	if err != nil {
		log.Printf("Error getting posts by user:%s, err:%s\n", u.ID, err.Error())
		utils.RespondWithInternalServerError(w)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, posts)
}
