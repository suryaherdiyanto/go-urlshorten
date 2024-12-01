package main

import (
	"context"
	"html/template"
	"log"
	"math/rand"
	"time"

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

func RandString(length int) string {
	charset := []byte("abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, length)

	for i := 0; i < length; i++ {
		rnd := rand.Intn(len(charset))
		b[i] = charset[int(rnd)]
	}

	return string(b)

}

func Add(i int, num int) int {
	return i + num
}

func main() {
	app := gin.Default()
	app.SetFuncMap(template.FuncMap{
		"add": Add,
	})

	store := cookie.NewStore([]byte("examplekey"))
	app.Use(sessions.Sessions("examplesession", store))
	db := database.New("root:root@tcp(127.0.0.1)/urlshorten")

	app.GET("/", func(ctx *gin.Context) {
		app.LoadHTMLFiles("views/master.tmpl", "views/menu.tmpl", "views/home.tmpl")
		session := sessions.Default(ctx)
		flash := session.Flashes()

		urls := []URL{}
		db.Select(&urls, "SELECT * FROM urls")
		ctx.HTML(200, "home.tmpl", gin.H{
			"urls":  &urls,
			"flash": flash,
		})
	})
	app.GET("create", func(ctx *gin.Context) {
		app.LoadHTMLFiles("views/master.tmpl", "views/menu.tmpl", "views/create.tmpl")
		session := sessions.Default(ctx)
		flash := session.Flashes()
		session.Save()

		ctx.HTML(200, "create.tmpl", gin.H{
			"flash": flash,
		})
	})
	app.POST("create", func(ctx *gin.Context) {
		err := ctx.Request.ParseForm()
		session := sessions.Default(ctx)

		if err != nil {
			log.Fatal(err)
		}

		url := ctx.Request.PostForm["url"]
		slug := RandString(6)

		con, cancle := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancle()

		_, err = db.ExecContext(con, "INSERT INTO urls(from_url,to_url) VALUES(?, ?)", "/r/"+slug, url[0])
		if err != nil {
			log.Fatalf("Something went wrong %v \n", err)
			ctx.Error(err)
		}
		session.AddFlash("http://localhost:8000/r/" + slug)
		session.Save()

		ctx.Redirect(302, "/create")
	})

	app.Run(":8000")
}
