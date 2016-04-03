package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func apiStreamJson(c *gin.Context) {
	client := <-SSE.NewClients
	ch := client.Channel
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return //log error & send some error
	}

	timeout := time.After(30 * time.Minute)

	for {
		select {
		case <-timeout:
			return
		case <-c.Done():
			return
		case event := <-ch:
			switch msg := event.(type) {
			case Tweet:
				c.SSEvent(MessageTweet, msg)
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

func initAPI(engine *gin.Engine) {

	engine.GET("/api/v1/stream.json", apiStreamJson)

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

}
