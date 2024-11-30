package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	app.LoadHTMLGlob("./views/*")

	app.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "home.tmpl", gin.H{})
	})

	app.Run(":8000")
}
