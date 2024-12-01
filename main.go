package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-urlshorten/database"
)

type URL struct {
	Id       int    `db:"id"`
	FromURL  string `db:"from_url"`
	ToURL    string `db:"to_url"`
	HitCount int    `db:"hit_count"`
}

type FormRequest struct {
	URL string `form:"url" binding:"required,url"`
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

func AppendError(bags map[string][]string, field string, message string) {
	if _, ok := bags[field]; !ok {
		bags[field] = []string{message}
	} else {
		bags[field] = append(bags[field], message)
	}
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
		session.Save()

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

		var data FormRequest
		if err = ctx.ShouldBind(&data); err != nil {
			errorBags := make(map[string][]string)

			for _, validationErr := range err.(validator.ValidationErrors) {
				var message string

				switch validationErr.Tag() {
				case "required":
					message = "The " + validationErr.Field() + " is required"
				case "url":
					message = "The " + validationErr.Field() + " must be a valid URL"
				default:
					message = validationErr.Error()
				}

				AppendError(errorBags, validationErr.Field(), message)
			}
			app.LoadHTMLFiles("views/master.tmpl", "views/menu.tmpl", "views/create.tmpl")

			ctx.HTML(400, "create.tmpl", gin.H{
				"errors": errorBags,
				"flash":  []string{},
			})
			return
		}

		slug := RandString(6)

		con, cancle := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancle()

		_, err = db.ExecContext(con, "INSERT INTO urls(from_url,to_url) VALUES(?, ?)", "/r/"+slug, data.URL)
		if err != nil {
			log.Fatalf("Something went wrong %v \n", err)
			ctx.Error(err)
		}
		session.AddFlash("http://localhost:3000/r/" + slug)
		session.Save()

		ctx.Redirect(302, "/create")
	})
	app.GET("/r/:slug", func(ctx *gin.Context) {
		var url URL

		con, cancle := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancle()

		err := db.GetContext(con, &url, "SELECT * FROM urls where from_url like ? LIMIT 1", "%"+ctx.Param("slug"))

		if err != nil {
			logEr := fmt.Sprintf("%v", err)
			fmt.Println(logEr)
			ctx.String(404, logEr)
		}

		_, err = db.ExecContext(con, "UPDATE urls set hit_count = ? where id = ?", url.HitCount+1, url.Id)

		if err != nil {
			logEr := fmt.Sprintf("%v", err)
			fmt.Println(logEr)
			ctx.String(500, logEr)
		}

		ctx.Redirect(302, url.ToURL)
	})
	app.POST("/delete/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		session := sessions.Default(ctx)

		con, cancle := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancle()

		_, err := db.ExecContext(con, "DELETE FROM urls WHERE id = ?", id)
		if err != nil {
			logEr := fmt.Sprintf("%v", err)
			fmt.Println(logEr)
			ctx.String(500, logEr)
		}

		session.AddFlash("Delete successfully!")
		session.Save()

		ctx.Redirect(302, "/")

	})

	app.Run(":8000")
}
