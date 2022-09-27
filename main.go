package main

import (
	"george0528/go-notion.git/controller"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	secret := os.Getenv("SESSION_SECRET")
	store := cookie.NewStore([]byte(secret))
	cookieOptions := sessions.Options{
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	store.Options(cookieOptions)
	router.Use(sessions.Sessions("origin_session", store))

	router.GET("/", func(ctx *gin.Context) {
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

	router.GET("/callback", func(ctx *gin.Context) {
		controller.Callback(ctx)
	})

	router.POST("/search", func(ctx *gin.Context) {
		controller.SearchNotion(ctx)
	})

	router.GET("/select/:id", func(ctx *gin.Context) {
		controller.Select(ctx)
	})

	router.POST("/schedule/:id", func(ctx *gin.Context) {
		controller.AddPages(ctx)
	})

	router.Run()
}