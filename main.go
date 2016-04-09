package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kerwindena/koma-bot/sse"
	"net/http"
	"os"
	"strings"
)

var Debug bool
var SSE sse.Provider

type SoundFS struct {
	fs http.FileSystem
}

func (s SoundFS) Open(name string) (http.File, error) {
	if strings.HasSuffix(name, ".wav") {
		return s.fs.Open(name)
	}
	return nil, os.ErrNotExist
}

func newSoundFS() SoundFS {
	s := SoundFS{http.Dir(getConfigString("sounds.dir"))}
	return s
}

func indexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{})
}

func main() {
	loadConfig()

	Debug = gin.IsDebugging()

	SSE = sse.NewProvider()
	go twitterListen()
	go processTweetSounds()

	router := gin.Default()
	router.StaticFS("/static", http.Dir("static"))
	router.StaticFS("/sounds", newSoundFS())
	router.LoadHTMLGlob("templates/*")

	initAPI(router)

	router.GET("/", indexPage)

	panic(router.Run("localhost:8000"))
}
