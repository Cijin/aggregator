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

type apiConfig struct {
	DB *database.Queries
}

func main() {
	apiCfg := apiConfig{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to load db: ", err)
	}

	dbQueries := database.New(db)
	apiCfg.DB = dbQueries

	port := os.Getenv("PORT")

	r := chi.NewRouter()
	r.Use(middleware.Cors())

	v1Router := chi.NewRouter()
	v1Router.Get("/readiness", v1.Readiness)
	v1Router.Get("/err", v1.Err)

	r.Mount("/v1", v1Router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Print("Listening on port: ", port)
	log.Fatal(server.ListenAndServe())
}
