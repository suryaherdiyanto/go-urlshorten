package main

import (
	"os"

	"github.com/go-urlshorten/app"
	"github.com/go-urlshorten/database"
	"github.com/go-urlshorten/handler"
)

type FormRequest struct {
	URL string `form:"url" binding:"required,url"`
}

type URL struct {
	Id       int    `db:"id"`
	FromURL  string `db:"from_url"`
	ToURL    string `db:"to_url"`
	HitCount int    `db:"hit_count"`
}

func main() {
	dbengine, ok := os.LookupEnv("DATABASE")
	if !ok {
		dbengine = "mysql"
	}

	dburl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		dburl = "root:@tcp(127.0.0.1)/"
	}

	app := app.NewApp(database.New(dbengine, dburl))
	app.Boot()
	h := handler.NewHandler(app.Db, app.Gin)
	h.SetupRouter()

	app.Run()

}
