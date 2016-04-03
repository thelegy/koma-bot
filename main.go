package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kerwindena/koma-bot/sse"
)

var Debug bool
var SSE sse.Provider

func main() {
	loadConfig()

	Debug = gin.IsDebugging()

	SSE = sse.NewProvider()
	go twitterListen()
	go processTweetSounds()

	router := gin.Default()

	initAPI(router)

	panic(router.Run("localhost:8000"))
}
