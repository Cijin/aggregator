package config

import "aggregator/internal/database"

type ApiConfig struct {
	Db *database.Queries
}

func NewConfig(db *database.Queries) *ApiConfig {
	return &ApiConfig{Db: db}
}
