package main

import (
	"aggregator/internal/database"
	"aggregator/middleware"
	v1 "aggregator/v1"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to load db: ", err)
	}

	defer db.Close()

	dbQueries := database.New(db)

	port := os.Getenv("PORT")

	r := chi.NewRouter()
	r.Use(middleware.Cors())

	v1Router := chi.NewRouter()
	v1Handler := v1.NewHandler(dbQueries)

	// start feed processing worker
	go v1Handler.FeedWorker()

	// v1 routes
	v1Router.Get("/readiness", v1Handler.Readiness)
	v1Router.Get("/err", v1Handler.Err)

	v1Router.Post("/users", v1Handler.CreateUser)

	v1Router.Get("/feeds", v1Handler.ListFeed)
	v1Router.Post("/feeds", v1Handler.MiddlewareAuth(v1Handler.CreateFeed))

	v1Router.Get("/feed_follows", v1Handler.MiddlewareAuth(v1Handler.GetFeed))
	v1Router.Post("/feed_follows", v1Handler.MiddlewareAuth(v1Handler.Follow))
	v1Router.Delete("/feed_follows/{feedFollowId}", v1Handler.MiddlewareAuth(v1Handler.Unfollow))

	v1Router.Get("/posts", v1Handler.MiddlewareAuth(v1Handler.ListPosts))

	r.Mount("/v1", v1Router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Print("Listening on port: ", port)
	log.Fatal(server.ListenAndServe())
}
