package main

import (
	"george0528/go-notion.git/controller"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

type Schedule struct {
	gorm.Model
  DbId string
  Name string
  FirstDay time.Time
  Interval int
	Time int
}

//DBマイグレート
func dbInit() {
	db, err := gorm.Open("sqlite3", "sqlite.db")
	if err != nil {
		panic("データベース開けず！（dbInit）")
	}
	db.AutoMigrate(&Schedule{})
	defer db.Close()
}

func main() {
	godotenv.Load(".env")
	dbInit()
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

	router.GET("/callback", func(ctx *gin.Context) {
		controller.Callback(ctx)
	})

	router.POST("/search", func(ctx *gin.Context) {
		controller.SearchNotion(ctx)
	})

	router.GET("/select/:id", func(ctx *gin.Context) {
		controller.Select(ctx)
	})

	router.Run()
}