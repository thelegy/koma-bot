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

func indexPage(v Version) func(*gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"Version": v,
		})
	}
}

func main() {
	config := loadConfig()
	version := getVersion()

	sse := sse.NewProvider()
	twitterApi := twitterConnect(config)

	config.ResolveUserIds(twitterApi)

	go twitterListen(twitterApi, config, sse.EventStream)

	go processTweetSounds(config, sse)

	loadRecentTweets(twitterApi, config)

	if !config.IsDebugging() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.StaticFS("/static", http.Dir("static"))
	router.StaticFS("/sounds", newSoundFS(config))
	router.LoadHTMLGlob("templates/*")

	initAPI(config, sse.NewClients, router)

	router.GET("/", indexPage(version))

	if len(os.Args) > 1 && os.Args[1] == "--docker" {
		panic(router.Run("0.0.0.0:8000"))
	} else {
		panic(router.Run(config.GetConfigString("address")))
	}
}
