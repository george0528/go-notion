package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

func Notion(c *gin.Context) {
	fmt.Println("test")
	baseUrl := "https://api.notion.com/v1/oauth/authorize";

	r, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	params := r.URL.Query()
	params.Add("client_id", os.Getenv("NOTION_CLIENT_ID"))
	params.Add("redirect_uri", "http://localhost:8080/callback")
	params.Add("response_type", "code")
	r.URL.RawQuery = params.Encode()
	redirectUrl := r.URL.String()
	fmt.Println(redirectUrl)

	c.Redirect(http.StatusMovedPermanently, redirectUrl)
}