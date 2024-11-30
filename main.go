package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-urlshorten/database"
)

type URL struct {
	Id       int    `db:"id"`
	FromURL  string `db:"from_url"`
	ToURL    string `db:"to_url"`
	HitCount int    `db:"hit_count"`
}

func main() {
	app := gin.Default()
	db := database.New("root:root@tcp(127.0.0.1)/urlshorten")
	app.LoadHTMLGlob("./views/*")

	app.GET("/", func(ctx *gin.Context) {
		urls := []URL{}
		db.Select(&urls, "SELECT * FROM urls")
		fmt.Printf("%v", urls)
		ctx.HTML(200, "home.tmpl", gin.H{
			"urls": &urls,
		})
	})

	app.Run(":8000")
}
