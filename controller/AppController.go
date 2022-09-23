package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// 定数
const notionUrl = "https://api.notion.com/v1"
var myToken string = ""

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

type RequestBody struct {
  Code string `json:"code"`
  GrantType string `json:"grant_type"`
  RedirectUri string `json:"redirect_uri"`
}

type TokenInfo struct {
	AccessToken   string `json:"access_token"`
	TokenType     string `json:"token_type"`
	BotID         string `json:"bot_id"`
	WorkspaceName string `json:"workspace_name"`
	WorkspaceIcon string `json:"workspace_icon"`
	WorkspaceID   string `json:"workspace_id"`
	Owner         struct {
		Type string `json:"type"`
		User struct {
			Object    string `json:"object"`
			ID        string `json:"id"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
			Type      string `json:"type"`
			Person    struct {
				Email string `json:"email"`
			} `json:"person"`
		} `json:" user"`
	} `json:"owner"`
}

func Callback(c *gin.Context) {
	baseUrl := "https://api.notion.com/v1/oauth/token"

	code := c.Query("code")

	requestBody := new(RequestBody)
	requestBody.Code = code
	requestBody.GrantType = "authorization_code"
	requestBody.RedirectUri = "http://localhost:8080/callback"

	// json
	jsonString, err := json.Marshal(requestBody)
  if err != nil {
    fmt.Println(err)
		return
  }

	// request作成
	request, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(jsonString))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 認証
	clinetId := os.Getenv("NOTION_CLIENT_ID")
	clientSecret :=  os.Getenv("NOTION_SECRET")
	request.SetBasicAuth(clinetId, clientSecret)

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
    Timeout: timeout,
	}
	request.Header.Set("Content-Type", "application/json")

	// clientで実行
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

	var tokenInfo TokenInfo
	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		fmt.Println(err)
		return;
	}

	fmt.Println(tokenInfo.AccessToken)
	myToken = tokenInfo.AccessToken
	c.HTML(200, "home.html", gin.H{})
}

type SearchRequestBody struct {
  Query string `json:"query"`
}
type SearchResponse struct {
	Object  string `json:"object"`
	Results []struct {
		Object      string      `json:"object"`
		ID          string      `json:"id"`
		Cover       interface{} `json:"cover"`
		Icon        interface{} `json:"icon"`
		CreatedTime time.Time   `json:"created_time"`
		CreatedBy   struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"created_by"`
		LastEditedBy struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"last_edited_by"`
		LastEditedTime time.Time `json:"last_edited_time"`
		Title          []struct {
			Type string `json:"type"`
			Text struct {
				Content string      `json:"content"`
				Link    interface{} `json:"link"`
			} `json:"text"`
			Annotations struct {
				Bold          bool   `json:"bold"`
				Italic        bool   `json:"italic"`
				Strikethrough bool   `json:"strikethrough"`
				Underline     bool   `json:"underline"`
				Code          bool   `json:"code"`
				Color         string `json:"color"`
			} `json:"annotations"`
			PlainText string      `json:"plain_text"`
			Href      interface{} `json:"href"`
		} `json:"title"`
		Description []interface{} `json:"description"`
		IsInline    bool          `json:"is_inline"`
		Properties  struct {
			NAMING_FAILED struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				Date struct {
				} `json:"date"`
			} `json:"日付"`
			NAMING_FAILED0 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Type        string `json:"type"`
				MultiSelect struct {
					Options []interface{} `json:"options"`
				} `json:"multi_select"`
			} `json:"タグ"`
			NAMING_FAILED1 struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Type  string `json:"type"`
				Title struct {
				} `json:"title"`
			} `json:"名前"`
		} `json:"properties"`
		Parent struct {
			Type      string `json:"type"`
			Workspace bool   `json:"workspace"`
		} `json:"parent"`
		URL      string `json:"url"`
		Archived bool   `json:"archived"`
	} `json:"results"`
	NextCursor     interface{} `json:"next_cursor"`
	HasMore        bool        `json:"has_more"`
	Type           string      `json:"type"`
	PageOrDatabase struct {
	} `json:"page_or_database"`
}

func SearchNotion(c *gin.Context) {
	fmt.Println("test")
	url := notionUrl + "/search"

	keyword := c.PostForm("keyword")
	requestBody := new(SearchRequestBody)
	requestBody.Query = keyword

	// json
	jsonString, err := json.Marshal(requestBody)
  if err != nil {
    fmt.Println(err)
		return
  }

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	if err != nil {
		fmt.Println(err)
		return
	}

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
    Timeout: timeout,
	}

	authorization := "Bearer "
	authorization += myToken
	r.Header.Set("Authorization", authorization)
	r.Header.Set("Notion-Version", "2022-06-28")
	r.Header.Set("Content-Type", "application/json")

	fmt.Println(r)
	request, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err)
		return;
	}

	var searchResponse SearchResponse
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		fmt.Println(err)
		return;
	}

	fmt.Println(searchResponse)
}