package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "this is AppController",
	})
}

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func Api(c *gin.Context) {
	url := "https://api.notion.com/v1/oauth/authorize"
	url = "https://jsonplaceholder.typicode.com/posts"
	r, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return;
	}

	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		fmt.Println(err)
		return;
	}

	fmt.Printf("%+v\n", posts)
}