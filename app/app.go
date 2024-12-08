package app

import (
	"html/template"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-urlshorten/database"
)

type App struct {
	Gin *gin.Engine
	Db  *database.Database
}

func Add(i int, num int) int {
	return i + num
}

func NewApp(db *database.Database) *App {

	return &App{
		Gin: gin.Default(),
		Db:  db,
	}
}

func (app *App) Run() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = ":8000"
	}

	app.Gin.Run(port)
}

func (app *App) Boot() {
	appkey, ok := os.LookupEnv("APP_KEY")
	if !ok {
		appkey = "somerandomkey"
	}

	app.Gin.SetFuncMap(template.FuncMap{
		"add": Add,
	})
	store := cookie.NewStore([]byte(appkey))
	app.Gin.Use(sessions.Sessions("_gin_session_", store))
}
