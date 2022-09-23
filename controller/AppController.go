package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// クエリパラメータ
	params := request.URL.Query()
	params.Add("userId","3")
	request.URL.RawQuery = params.Encode()

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
    Timeout: timeout,
	}

	r, err := client.Do(request)
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