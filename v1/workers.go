package v1

import (
	"aggregator/internal/database"
	"aggregator/utils"
	"context"
	"database/sql"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const UniqueViolationCode = "23505"

var offset int32 = 0

func (v *v1) FeedWorker() {
	for range time.Tick(time.Second * 10) {
		v.ProcessFeed(10)
	}
}

func (v *v1) ProcessFeed(n int32) {
	// get n feeds from database
	feedFetchParams := database.ListFeedsToFetchParams{
		Limit:  n,
		Offset: offset,
	}
	feeds, err := v.Db.ListFeedsToFetch(context.Background(), feedFetchParams)
	if err != nil {
		log.Fatal("Error fetching feeds: ", err)
	}

	var wg sync.WaitGroup

	for _, feed := range feeds {
		wg.Add(1)
		go v.FetchAndMarkFeed(feed, &wg)
	}

	wg.Wait()

	offset += n
}

func (v *v1) FetchAndMarkFeed(feed database.Feed, wg *sync.WaitGroup) {
	defer wg.Done()

	u, err := url.Parse(feed.Url)
	if err != nil {
		log.Fatal("Url is not valid: ", err)
		return
	}

	xml, err := utils.FetchFeed(u)
	if err != nil {
		log.Fatal("Err fetching feed: ", err)
		return
	}

	// mark feed fetched
	err = v.Db.MarkFeedFetched(
		context.Background(),
		database.MarkFeedFetchedParams{
			ID:            feed.ID,
			LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedAt:     time.Now(),
		},
	)
	if err != nil {
		log.Print("Error marking feed as fetched: ", err)
	}

	for _, post := range xml.Channel.Item {
		v.CreatePost(&post, feed.ID)
	}
}

func (v *v1) CreatePost(post *utils.Item, feedId uuid.UUID) {
	var publishedAt sql.NullTime
	publishedAt.Valid = true

	lastPublishedAt, err := time.Parse(time.RFC1123Z, post.PubDate)
	if err != nil {
		publishedAt.Valid = false
		log.Printf("Error parsing feed: %s date: %s\n", feedId, post.PubDate)
	}

	publishedAt.Time = lastPublishedAt

	var description sql.NullString
	description.Valid = true
	description.String = post.Description

	if post.Description == "" {
		description.Valid = false
	}

	postParam := database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       post.Title,
		Url:         post.Link,
		Description: description,
		PublishedAt: publishedAt,
		FeedID:      feedId,
	}

	_, err = v.Db.CreatePost(context.Background(), postParam)
	if e, ok := err.(*pq.Error); ok {
		if e.Code.Class() != UniqueViolationCode {
			log.Printf("%s already exists, err:%s", post.Link, err.Error())
		}
	}
}
