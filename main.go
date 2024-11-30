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

func Flash(ctx *gin.Context, key string, value interface{}) (string, bool) {
	session := sessions.Default(ctx)

	if value == nil {
		session.AddFlash(value, key)
	}

	message := session.Flashes(key)

	if message == nil {
		return "", false
	}

	return message[0].(string), true
}

func main() {
	app := gin.Default()
	app.SetFuncMap(template.FuncMap{
		"Flash": Flash,
	})

	store := cookie.NewStore([]byte("examplekey"))
	app.Use(sessions.Sessions("examplesession", store))
	db := database.New("root:root@tcp(127.0.0.1)/urlshorten")
	app.LoadHTMLGlob("./views/*")

	app.GET("/", func(ctx *gin.Context) {
		urls := []URL{}
		db.Select(&urls, "SELECT * FROM urls")
		ctx.HTML(200, "home.tmpl", gin.H{
			"urls": &urls,
		})
	})

	app.Run(":8000")
}
