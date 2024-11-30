package database

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func New(conn string) *sqlx.DB {
	db, err := sqlx.Connect("mysql", conn)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	return db
}
