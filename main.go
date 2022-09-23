package main

import (
	"george0528/go-notion.git/controller"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", func(ctx *gin.Context){
		ctx.HTML(200, "index.html", gin.H{})
	})

	router.GET("/index", func(ctx *gin.Context) {
		controller.Index(ctx)
	})

	router.GET("/api", func(ctx *gin.Context) {
		controller.Api(ctx)
	})

	router.GET("/notion", func(ctx *gin.Context) {
		controller.Notion(ctx)
	})

	router.Run()
}