package main

import (
	"george0528/go-notion.git/controller"

	"github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    router.LoadHTMLGlob("templates/*.html")

    router.GET("/", func(ctx *gin.Context){
        ctx.HTML(200, "index.html", gin.H{})
    })

	router.GET("/index", func(ctx *gin.Context) {
		controller.Index(ctx)
	})

    router.Run()
}