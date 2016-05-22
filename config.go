package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/chimeracoder/anaconda"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	configViper *viper.Viper
	soundViper  *viper.Viper
	soundMutex  sync.RWMutex
	soundMap    map[string]*Sound
	StreamInfo  [2]*tweetStreamInfo
	userIds     []int64
}

func (conf *Config) IsDebugging() bool {
	return conf.configViper.GetBool("debug")
}

func (conf *Config) loadStreamInfo() {
	streamInfo1 := newTweetStreamInfo()
	streamInfo2 := newTweetStreamInfo()
	conf.configViper.UnmarshalKey("twitter.stream1", streamInfo1)
	conf.configViper.UnmarshalKey("twitter.stream2", streamInfo2)
	conf.StreamInfo[0] = streamInfo1
	conf.StreamInfo[1] = streamInfo2
}

func (conf *Config) ResolveUserIds(api *anaconda.TwitterApi) {
	if len(conf.userIds) > 0 {
		return
	}
	var userNames bytes.Buffer
	for _, tsi := range conf.StreamInfo {
		for _, user := range tsi.Users {
			userNames.WriteString(",")
			userNames.WriteString(user)
		}
	}
	userNames.Next(1)
	val := url.Values{}
	users, err := api.GetUsersLookup(userNames.String(), val)
	if err != nil {
		//log error
		return
	}
	for _, user := range users {
		conf.userIds = append(conf.userIds, user.Id)
	}
}

func (conf *Config) GetTweetFilter() (string, string) {
	var track bytes.Buffer
	var follow bytes.Buffer
	for _, tsi := range conf.StreamInfo {
		for _, hashtag := range tsi.Hashtags {
			track.WriteString(",#")
			track.WriteString(hashtag)
		}
	}
	track.Next(1)
	for _, uid := range conf.userIds {
		follow.WriteString(",")
		follow.WriteString(strconv.FormatInt(uid, 10))
	}
	follow.Next(1)
	return track.String(), follow.String()
}

func (conf *Config) StoreTweet(t Tweet) {
	for _, tsi := range conf.StreamInfo {
		if tsi.ContainsTweet(t) {
			tsi.Add(&t)
		}
	}
}

func (conf *Config) iterateSounds() <-chan *Sound {
	sounds := make(chan *Sound)
	go func(c chan<- *Sound) {
		conf.soundMutex.RLock()
		defer conf.soundMutex.RUnlock()
		defer close(c)

		timeout := time.After(100 * time.Millisecond)

		for _, sound := range conf.soundMap {
			select {
			case c <- sound:
			case <-timeout:
				return
			}
		}

	}(sounds)
	return sounds
}

func (conf *Config) updateSounds() {
	sounds := make(map[string]*Sound)
	for _, key := range conf.soundViper.AllKeys() {
		sound := &Sound{}
		err := conf.soundViper.UnmarshalKey(key, sound)
		if err != nil {
			fmt.Println("Sound error:", err) //error handling, no panic
			continue
		}
		sound.Name = key
		err = sound.CompileRegexpr()
		if err != nil {
			fmt.Println("Sound error:", err) //error handling, no panic
			continue
		}
		sounds[key] = sound
	}

	conf.soundMutex.Lock()
	defer conf.soundMutex.Unlock()

	conf.soundMap = sounds
}

func soundConfigChanged(conf *Config) func(fsnotify.Event) {
	return func(_ fsnotify.Event) {
		conf.updateSounds()
	}
}

func loadConfig() *Config {
	conf := &Config{
		configViper: viper.New(),
		soundViper:  viper.New(),
	}

	conf.configViper.SetConfigName("koma_bot")
	conf.configViper.AddConfigPath("/etc/koma_bot/")
	conf.configViper.AddConfigPath(".")

	err := conf.configViper.ReadInConfig() // Find and read the config file
	if err != nil {                        // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	conf.configViper.SetDefault("debug", false)

	conf.loadStreamInfo()

	conf.soundViper.SetConfigName("koma_bot_sounds")
	conf.soundViper.AddConfigPath("/etc/koma_bot/")
	conf.soundViper.AddConfigPath(".")

	err = conf.soundViper.ReadInConfig() // Find and read the config file
	if err != nil {                      // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	conf.soundViper.OnConfigChange(soundConfigChanged(conf))
	conf.soundViper.WatchConfig()

	conf.updateSounds()

	return conf
}

func (conf *Config) GetConfigString(name string) string {
	return conf.configViper.GetString(name)
}
