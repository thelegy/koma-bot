package main

import (
	"regexp"
	"strings"
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
	pos := s.regex.FindAllStringIndex(t.Text, -1)
	if len(pos) == 0 {
		return nil
	}
	return pos
}

func tweetSounds(t Tweet) {
	positions := make([]*Sound, len(t.Text)+1)
	for sound := range iterateSounds() {
		pos := sound.GetPosition(t)
		for _, p := range pos {
			positions[p[0]] = sound
		}
	}
	for _, sound := range positions {
		if sound == nil {
			continue
		}
		SSE.EventStream <- sound
	}
}

func processTweetSounds() {
	for {
		c := <-SSE.NewClients
		for m := range c.Channel {
			switch msg := m.(type) {
			default:
			case Tweet:
				tweetSounds(msg)
			}
		}
	}
}
