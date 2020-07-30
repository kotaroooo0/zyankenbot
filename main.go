package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	r.GET("/twitter_webhook", HandlerCrcCheck)
	r.POST("/twitter_webhook", HandlerTwitterActivity)
	r.Run(":3000")
}

func HandlerCrcCheck(c *gin.Context) {
	// Requestを受ける
	req := GetCrcCheckRequest{}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	// CrcTokenを生成し、Responseに詰める
	mac := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	mac.Write([]byte(req.CrcToken))
	res := GetCrcCheckResponse{
		Token: "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil)),
	}
	// Responseを返す
	c.JSON(http.StatusOK, res)
}

type GetCrcCheckRequest struct {
	CrcToken string `json:"crc_token" form:"crc_token" binding:"required"`
}

type GetCrcCheckResponse struct {
	Token string `json:"response_token"`
}

func HandlerTwitterActivity(c *gin.Context) {
	// Requestを受ける
	req := PostTwitterActivityRequest{}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// リプライがない、もしくはユーザが不正な場合は400を返す
	if len(req.TweetCreateEvents) < 1 || req.UserID == req.TweetCreateEvents[0].User.IDStr {
		c.JSON(http.StatusBadRequest, "no reply or invalid user")
		return
	}

	// リプライの内容を取得
	replyText := req.TweetCreateEvents[0].Text

	// @snowfall_botを消す
	replyText = strings.Replace(replyText, "@snowfall_bot ", "", -1)

	// 返信内容を生成
	content := ContentByReplyText(replyText)

	// リプライを返す
	params := url.Values{}
	params.Set("in_reply_to_status_id", req.TweetCreateEvents[0].TweetIDStr)
	twitterApiClient := NewTwitterApiClient()
	status, err := twitterApiClient.PostTweet(fmt.Sprintf("@%s %s", req.TweetCreateEvents[0].User.ScreenName, content), params)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func NewTwitterApiClient() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET"))
	return api
}

func ContentByReplyText(replyText string) string {
	if replyText == "ぐー" || replyText == "ちょき" || replyText == "ぱー" {
		rand.Seed(time.Now().UnixNano())
		if rand.Float32() > 0.3 {
			return "俺の勝ち！何で負けたか、明日まで考えといてください"
		} else {
			return "今日は負けを認めます。ただ、勝ち逃げは許しませんよ"
		}
	}
	return "ぐー、ちょき、ぱーのいずれかで"
}

type PostTwitterActivityRequest struct {
	UserID            string             `json:"for_user_id" form:"for_user_id" binding:"required"`
	TweetCreateEvents []TweetCreateEvent `json:"tweet_create_events" form:"tweet_create_events" binding:"required"`
}

type TweetCreateEvent struct {
	TweetIDStr string `json:"id_str" form:"id_str" binding:"required"`
	Text       string `json:"text" form:"text" binding:"required"`
	User       struct {
		IDStr      string `json:"id_str" form:"id_str" binding:"required"`
		ScreenName string `json:"screen_name" form:"screen_name" binding:"required"`
	} `json:"user" form:"user" binding:"required"`
}
