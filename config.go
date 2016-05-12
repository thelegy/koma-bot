package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	configViper *viper.Viper
	soundViper  *viper.Viper
	soundMutex  sync.RWMutex
	soundMap    map[string]*Sound
}

func (conf *Config) IsDebugging() bool {
	return conf.configViper.GetBool("debug")
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
