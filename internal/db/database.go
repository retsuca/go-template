package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-template/internal/config"
	log "go.uber.org/zap"
)

var DB *sql.DB

func init() {

	queryString := fmt.Sprintf("host=%s port=5432 dbname=%s user=%s  password=%s sslmode=disable", config.Get(config.DB_ADDRESS), config.Get(config.DB_NAME), config.Get(config.DB_USERNAME), config.Get(config.DB_PASSWORD))
	db, err := sql.Open("postgres", queryString)

	if err != nil {
		log.S().Panic("error when try to connect DB", err)
	}

	DB = db
}
