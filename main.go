package main

import (
	"text/template"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-urlshorten/database"
)

type URL struct {
	Id       int    `db:"id"`
	FromURL  string `db:"from_url"`
	ToURL    string `db:"to_url"`
	HitCount int    `db:"hit_count"`
}

func Flash(ctx *gin.Context, key string, value interface{}) interface{} {
	session := sessions.Default(ctx)

	if value == nil {
		session.AddFlash(value, key)
	}

	message := session.Flashes(key)

	if message == nil {
		return nil
	}

	return message[0].(string)
}

func main() {
	app := gin.Default()
	app.SetFuncMap(template.FuncMap{
		"Flash": Flash,
	})

	store := cookie.NewStore([]byte("examplekey"))
	app.Use(sessions.Sessions("examplesession", store))
	db := database.New("root:root@tcp(127.0.0.1)/urlshorten")

	app.GET("/", func(ctx *gin.Context) {
		app.LoadHTMLFiles("views/master.tmpl", "views/menu.tmpl", "views/home.tmpl")
		urls := []URL{}
		db.Select(&urls, "SELECT * FROM urls")
		ctx.HTML(200, "home.tmpl", gin.H{
			"urls": &urls,
		})
	})
	app.GET("create", func(ctx *gin.Context) {
		app.LoadHTMLFiles("views/master.tmpl", "views/menu.tmpl", "views/create.tmpl")
		ctx.HTML(200, "create.tmpl", gin.H{})
	})

	app.Run(":8000")
}
