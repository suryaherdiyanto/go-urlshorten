package main

import (
	"github.com/go-urlshorten/app"
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
	app := app.NewApp()
	app.Boot()
	gi := app.Gin

	h := handler.NewHandler(app)

	gi.GET("/", h.Home)
	gi.GET("create", h.Create)
	gi.POST("create", h.Store)
	gi.GET("/r/:slug", h.Redirect)
	gi.POST("/delete/:id", h.Delete)

	app.Run()

}
