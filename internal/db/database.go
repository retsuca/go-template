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

	queryString := fmt.Sprintf("host=localhost port=5432 user=%s dbname=postgres password=%s sslmode=disable", config.Get("DBUser"), config.Get("DBPW"))
	db, err := sql.Open("postgres", queryString)

	if err != nil {
		log.S().Panic("error when try to connect DB", err)
	}

	DB = db
}
