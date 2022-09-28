package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 定数
const notionUrl = "https://api.notion.com/v1"

// 関数
func getClient() *http.Client {
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	return client
}

func setHeaders(r *http.Request, c *gin.Context) *http.Request {
	session := sessions.Default(c);
	authorization := "Bearer "
	myToken := session.Get("token")
	if myToken == nil {
		myToken = ""
		// todo エラーを吐くようにする
	}
	authorization += myToken.(string)
	r.Header.Set("Authorization", authorization)
	r.Header.Set("Notion-Version", "2022-06-28")
	r.Header.Set("Content-Type", "application/json")
	return r
}

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
	url := "https://jsonplaceholder.typicode.com/posts"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// クエリパラメータ
	params := request.URL.Query()
	params.Add("userId", "3")
	request.URL.RawQuery = params.Encode()

	client := getClient()
	r, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", posts)
}

func Notion(c *gin.Context) {
	fmt.Println("test")
	baseUrl := "https://api.notion.com/v1/oauth/authorize"

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
	Code        string `json:"code"`
	GrantType   string `json:"grant_type"`
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
	clientSecret := os.Getenv("NOTION_SECRET")
	request.SetBasicAuth(clinetId, clientSecret)

	client := getClient()
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
		return
	}

	var tokenInfo TokenInfo
	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(tokenInfo.AccessToken)
	myToken := tokenInfo.AccessToken
	session := sessions.Default(c)
	session.Set("token", myToken)
	session.Save()
	c.HTML(200, "home.html", gin.H{})
}

type SearchRequestBody struct {
	Query string `json:"query"`
}

type Result struct {
	Object string      `json:"object"`
	ID     string      `json:"id"`
	Cover  interface{} `json:"cover"`
	Icon   struct {
		Type  string `json:"type"`
		Emoji string `json:"emoji"`
	} `json:"icon"`
	CreatedTime time.Time `json:"created_time"`
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
	Properties  map[string]map[string]interface{} `json:"properties"`
	Parent struct {
		Type      string `json:"type"`
		Workspace bool   `json:"workspace"`
	} `json:"parent"`
	URL      string `json:"url"`
	Archived bool   `json:"archived"`
}
type SearchResponse struct {
	Object         string      `json:"object"`
	Results        []Result    `json:"results"`
	NextCursor     interface{} `json:"next_cursor"`
	HasMore        bool        `json:"has_more"`
	Type           string      `json:"type"`
	PageOrDatabase struct {
	} `json:"page_or_database"`
}

type Title struct {
	Text string
	Id   string
}

func SearchNotion(c *gin.Context) {
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

	client := getClient()
	r = setHeaders(r, c)

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
		return
	}

	fmt.Println(string(body))

	var searchResponse SearchResponse
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("-------------------")

	var databases []Result
	for _, v := range searchResponse.Results {
		if v.Object == "database" {
			databases = append(databases, v)
		}
	}

	var titles []Title
	for _, v := range databases {
		title := Title{
			Text: v.Title[0].Text.Content,
			Id:   v.ID,
		}
		titles = append(titles, title)
	}

	c.HTML(200, "home.html", gin.H{
		"titles": titles,
	})
}

func Select(c *gin.Context) {
	id := c.Param("id")
	url := notionUrl + "/databases/" + id

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := getClient()
	r = setHeaders(r, c)

	request, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))

	var selectResponse Result
	if err := json.Unmarshal(body, &selectResponse); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("----------------")
	var properties []string

	for k, v := range selectResponse.Properties {
		if v["type"] == "date" {
			properties = append(properties, k)
		}
	}

	c.HTML(200, "select.html", gin.H{
		"id": id,
		"properties": properties,
	})
}

type AddPageRequest struct {
	Parent struct {
		DatabaseID string `json:"database_id"`
	} `json:"parent"`
	Properties map[string]Property `json:"properties"`
}
type Property map[string]interface{}
type NotionText struct {
	Type string `json:"type"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

func addPageRequest(requestBody AddPageRequest, url string, c *gin.Context) *http.Response {
	// json
	jsonString, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(r)

	client := getClient()
	r = setHeaders(r, c)
	request, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println(string(body))

	return request
}

func createRequestBody(id string, name string, dateName string, day string) *AddPageRequest {
	requestBody := new(AddPageRequest)
	requestBody.Parent.DatabaseID = id
	requestBody.Properties = make(map[string]Property)
	requestBody.Properties[dateName] = make(Property)
	requestBody.Properties[dateName]["date"] = map[string]string{"start": day, "time_zone": "Asia/Tokyo"}
	var notionText interface{} = NotionText{
		Type: "text",
		Text: struct{Content string "json:\"content\""}{Content: name},
	}
	requestBody.Properties["名前"] = make(Property)
	requestBody.Properties["名前"]["title"] = []interface{}{notionText}
	return requestBody
}

const layout = "2006-01-02T15:04"
func timeToString(t time.Time) string {
	str := t.Format(layout)
	return str
}

func stringToTime(str string) time.Time {
	t, _ := time.Parse(layout, str)
	return t
}


func AddPages(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	firstDay := c.PostForm("firstDay")
	dateName := c.PostForm("dateName")
	intervalStr := c.PostForm("interval")
	numStr := c.PostForm("num")

	url := notionUrl + "/pages"

	num, _ := strconv.Atoi(numStr)
	interval, _ := strconv.Atoi(intervalStr)
	for i := 0; i < num; i++ {
		date := firstDay
		dateTime := stringToTime(date)
		if i != 0 {
			dateTime = dateTime.AddDate(0, 0, interval * i)
		}

		dateStr := timeToString(dateTime)
		requestBody := createRequestBody(id, name, dateName, dateStr)
		addPageRequest(*requestBody, url, c)
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}
