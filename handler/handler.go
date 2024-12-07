package handler

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-urlshorten/app"
)

type Handler struct {
	App *app.App
}

type URL struct {
	Id       int    `db:"id"`
	FromURL  string `db:"from_url"`
	ToURL    string `db:"to_url"`
	HitCount int    `db:"hit_count"`
}

func NewHandler(app *app.App) *Handler {
	return &Handler{App: app}
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

func AppendError(bags map[string][]string, field string, message string) {
	if _, ok := bags[field]; !ok {
		bags[field] = []string{message}
	} else {
		bags[field] = append(bags[field], message)
	}
}

func (h *Handler) Home(ctx *gin.Context) {
	h.App.Gin.LoadHTMLFiles("views/master.tmpl", "views/menu.tmpl", "views/home.tmpl")
	session := sessions.Default(ctx)
	flash := session.Flashes()
	session.Save()

	urls := []URL{}
	err, _ := h.App.Db.All(&urls, "urls", "*")
	if err != nil {
		ctx.String(500, fmt.Sprintf("Something went wrong: %v", err))
		return
	}

	ctx.HTML(200, "home.tmpl", gin.H{
		"urls":  &urls,
		"flash": flash,
	})
}

func (h *Handler) Create(ctx *gin.Context) {
	h.App.Gin.LoadHTMLFiles("views/master.tmpl", "views/menu.tmpl", "views/create.tmpl")
	session := sessions.Default(ctx)
	flash := session.Flashes()
	session.Save()

	ctx.HTML(200, "create.tmpl", gin.H{
		"flash": flash,
	})
}

func (h *Handler) Store(ctx *gin.Context) {
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
		h.App.Gin.LoadHTMLFiles("views/create.tmpl")

		ctx.HTML(400, "create.tmpl", gin.H{
			"errors": errorBags,
			"flash":  []string{},
		})
		return
	}

	slug := RandString(6)

	err, _ = h.App.Db.Create(map[string]interface{}{"from_url": "r/" + slug, "to_url": data.URL}, "urls")
	if err != nil {
		log.Fatalf("Something went wrong %v \n", err)
		ctx.Error(err)
		ctx.String(500, "Something went wrong %v \n", err)
	}
	session.AddFlash("http://localhost:3000/r/" + slug)
	session.Save()

	ctx.Redirect(302, "/create")
}

func (h *Handler) Redirect(ctx *gin.Context) {
	var url URL

	err := h.App.Db.GetRaw("SELECT * FROM urls where from_url like ? LIMIT 1", &url, "%"+ctx.Param("slug"))

	if err != nil {
		logEr := fmt.Sprintf("%v", err)
		fmt.Println(logEr)
		ctx.String(404, fmt.Sprintf("Could not find that slug: %s", ctx.Param("slug")))
		return
	}

	err, _ = h.App.Db.Update(map[string]interface{}{"hit_count": url.HitCount + 1}, "urls", url.Id)

	if err != nil {
		logEr := fmt.Sprintf("%v", err)
		fmt.Println(logEr)
		ctx.String(500, logEr)
	}

	ctx.Redirect(302, url.ToURL)
}

func (h *Handler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	session := sessions.Default(ctx)

	if err != nil {
		ctx.String(404, "Resource not found")
	}

	err, _ = h.App.Db.Delete("urls", id)
	if err != nil {
		logEr := fmt.Sprintf("%v", err)
		fmt.Println(logEr)
		ctx.String(500, logEr)
	}

	session.AddFlash("Delete successfully!")
	session.Save()

	ctx.Redirect(302, ctx.Request.Referer())
}
