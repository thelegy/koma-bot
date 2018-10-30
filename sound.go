package main

import (
	"regexp"
	"strings"

	"github.com/thelegy/koma-bot/sse"
)

type Sound struct {
	Name      string
	Regex     []string
	Retweeted bool
	regex     *regexp.Regexp
}

func (s *Sound) CompileRegexpr() error {
	reg := "(?i)" + strings.Join(s.Regex, "|")
	regex, err := regexp.Compile(reg)
	if err != nil {
		return err
	}
	s.regex = regex
	return nil
}

func (s *Sound) GetPosition(t Tweet) [][]int {
	if s.Retweeted && t.RetweetedStatus == nil {
		return nil
	}
	pos := s.regex.FindAllStringIndex(t.FullText, -1)
	if len(pos) == 0 {
		return nil
	}
	return pos
}

func tweetSounds(conf *Config, sse sse.Provider, t Tweet) {
	positions := make([]*Sound, len(t.FullText)+1)
	for sound := range conf.iterateSounds() {
		pos := sound.GetPosition(t)
		for _, p := range pos {
			positions[p[0]] = sound
		}
	}
	for _, sound := range positions {
		if sound == nil {
			continue
		}
		sse.EventStream <- sound
	}
}

func processTweetSounds(conf *Config, sse sse.Provider) {
	for {
		c := <-sse.NewClients
		for m := range c.Channel {
			switch msg := m.(type) {
			default:
			case Tweet:
				tweetSounds(conf, sse, msg)
			}
		}
	}
}
