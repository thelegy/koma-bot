package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kerwindena/koma-bot/sse"
)

func apiStreamJson(clients <-chan sse.Client, ts *TweetStorage) func(c *gin.Context) {
	return func(c *gin.Context) {
		client := <-clients
		ch := client.Channel
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			return //log error & send some error
		}

		timeout := time.After(30 * time.Minute)

		for _, t := range ts.getTweets() {
			if t != nil {
				c.SSEvent(MessageTweet1, t)
				flusher.Flush()
			}
		}

		for {
			select {
			case <-timeout:
				return
			case <-c.Done():
				return
			case event := <-ch:
				switch msg := event.(type) {
				case Tweet:
					c.SSEvent(MessageTweet1, msg)
					flusher.Flush()
				case *Sound:
					c.SSEvent(MessageSound, msg.Name)
					flusher.Flush()
				default:
					continue
				}
			}
		}
	}
}

func initAPI(clients <-chan sse.Client, engine *gin.Engine, tweetStorage *TweetStorage) {

	engine.GET("/api/v1/stream.json", apiStreamJson(clients, tweetStorage))

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

}
