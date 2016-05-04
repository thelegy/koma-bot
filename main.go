package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kerwindena/koma-bot/sse"

	"net/http"
	"os"
	"strings"
)

type SoundFS struct {
	fs http.FileSystem
}

func (s SoundFS) Open(name string) (http.File, error) {
	if strings.HasSuffix(name, ".wav") {
		return s.fs.Open(name)
	}
	return nil, os.ErrNotExist
}

func newSoundFS(c *Config) SoundFS {
	s := SoundFS{http.Dir(c.GetConfigString("sounds.dir"))}
	return s
}

func indexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{})
}

func main() {
	config := loadConfig()

	sse := sse.NewProvider()
	go twitterListen(config, sse.EventStream)
	go processTweetSounds(config, sse)

	router := gin.Default()
	router.StaticFS("/static", http.Dir("static"))
	router.StaticFS("/sounds", newSoundFS(config))
	router.LoadHTMLGlob("templates/*")

	initAPI(sse.NewClients, router)

	router.GET("/", indexPage)

	panic(router.Run("localhost:8000"))
}
