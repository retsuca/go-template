package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Register postgres driver
	"go-template/internal/config"
	"go-template/pkg/logger"
	"go.uber.org/zap"
)

var DB *sql.DB

func init() {
	dataSource := fmt.Sprintf(
		"host=%s port=5432 dbname=%s user=%s  password=%s sslmode=disable",
		config.Get(config.DB_ADDRESS),
		config.Get(config.DB_NAME),
		config.Get(config.DB_USER),
		config.Get(config.DB_PW),
	)

	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	DB = db
}
