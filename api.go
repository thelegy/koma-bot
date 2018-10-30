package main

import (
	"net/http"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thelegy/koma-bot/sse"
)

func apiStreamJson(conf *Config, clients <-chan sse.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		client := <-clients
		ch := client.Channel
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			return //log error & send some error
		}

		timeout := time.After(30 * time.Minute)

		for i, streamInfo := range conf.StreamInfo {
			for _, t := range streamInfo.getTweets() {
				if t != nil {
					c.SSEvent(MessageTweet+strconv.Itoa(i+1), t)
				}
			}
		}

		c.SSEvent("", "")
		flusher.Flush()

		for {
			keepAliveTimer := time.After(55 * time.Second)
			select {
			case <-timeout:
				return
			case <-c.Done():
				return
			case <-keepAliveTimer:
				c.Writer.WriteString(":\n\n")
				flusher.Flush()
			case event := <-ch:
				switch msg := event.(type) {
				case Tweet:
					messageSent := false
					for i, streamInfo := range conf.StreamInfo {
						if streamInfo.ContainsTweet(msg) {
							c.SSEvent(MessageTweet+strconv.Itoa(i+1), msg)
							messageSent = true
						}
					}
					if !messageSent {
						c.Writer.WriteString(":\n\n")
					}
					flusher.Flush()
				case syscall.Signal:
					if msg == syscall.SIGUSR1 {
						c.SSEvent(MessageReload, "")
					} else {
						c.Writer.WriteString(":\n\n")
					}
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

func initAPI(conf *Config, clients <-chan sse.Client, engine *gin.Engine) {

	engine.GET("/api/v1/stream.json", apiStreamJson(conf, clients))

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

}
