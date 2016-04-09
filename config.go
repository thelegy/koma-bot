package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/fsnotify.v1"
	"sync"
	"time"
)

var (
	configViper *viper.Viper
	soundViper  *viper.Viper
	soundMutex  = sync.RWMutex{}
	soundMap    map[string]*Sound
)

func iterateSounds() <-chan *Sound {
	sounds := make(chan *Sound)
	go func(c chan<- *Sound) {
		soundMutex.RLock()
		defer soundMutex.RUnlock()
		defer close(c)

		timeout := time.After(100 * time.Millisecond)

		for _, sound := range soundMap {
			select {
			case c <- sound:
			case <-timeout:
				return
			}
		}

	}(sounds)
	return sounds
}

func updateSounds() {
	sounds := make(map[string]*Sound)
	for _, key := range soundViper.AllKeys() {
		sound := &Sound{}
		err := soundViper.UnmarshalKey(key, sound)
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

	soundMutex.Lock()
	defer soundMutex.Unlock()

	soundMap = sounds
}

func soundConfigChanged(event fsnotify.Event) {
	updateSounds()
}

func loadConfig() {
	configViper = viper.New()
	configViper.SetConfigName("koma_bot")
	configViper.AddConfigPath("/etc/koma_bot/")
	configViper.AddConfigPath(".")
	err := configViper.ReadInConfig() // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	soundViper = viper.New()
	soundViper.SetConfigName("koma_bot_sounds")
	soundViper.AddConfigPath("/etc/koma_bot/")
	soundViper.AddConfigPath(".")
	err = soundViper.ReadInConfig() // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	soundViper.OnConfigChange(soundConfigChanged)
	soundViper.WatchConfig()
	updateSounds()
}

func getConfigString(name string) string {
	return configViper.GetString(name)
}
