package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-template/internal/config"
)

var DB *sql.DB

func init() {

	queryString := fmt.Sprintf("host=localhost port=5432 user=%s dbname=postgres password=%s sslmode=disable", config.DBUser, config.DBPW)
	db, err := sql.Open("postgres", queryString)

	if err != nil {
		panic(err)
	}

	DB = db

}
