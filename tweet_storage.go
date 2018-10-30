package main

import (
	"bytes"
	"errors"
	"strings"
	"sync"

	"github.com/thelegy/koma-bot/sse"
)

const TweetStorageMaxCapacity = 100

type tweetStreamInfo struct {
	Hashtags []string
	Users    []string
	tweets   [TweetStorageMaxCapacity]*Tweet
	mutex    *sync.RWMutex
}

func newTweetStreamInfo() *tweetStreamInfo {
	tsi := &tweetStreamInfo{
		mutex: &sync.RWMutex{},
	}
	return tsi
}

func (tsi *tweetStreamInfo) GetTweetFilter() string {
	var filter bytes.Buffer
	for _, hashtag := range tsi.Hashtags {
		filter.WriteString(" OR #")
		filter.WriteString(hashtag)
	}
	for _, user := range tsi.Users {
		filter.WriteString(" OR from:")
		filter.WriteString(user)
	}
	filter.Next(4)
	return filter.String()
}

func (tsi *tweetStreamInfo) Add(t *Tweet) error {
	j := -1

	tsi.mutex.Lock()
	defer tsi.mutex.Unlock()

	for j < TweetStorageMaxCapacity-1 {
		if tsi.tweets[j+1] == nil {
			j++
			continue
		}

		if t.Id == tsi.tweets[j+1].Id {
			return errors.New("Tweet already exists")
		}

		if t.Id < tsi.tweets[j+1].Id {
			break
		}
		j++
	}

	if j < 0 {
		return nil
	}

	if j+1 > TweetStorageMaxCapacity {
		j = TweetStorageMaxCapacity - 1
	}

	for i := 0; i < j; i++ {
		tsi.tweets[i] = tsi.tweets[i+1]
	}
	tsi.tweets[j] = t

	return nil
}

func (tsi *tweetStreamInfo) getTweets() []*Tweet {
	tweets := make([]*Tweet, TweetStorageMaxCapacity)

	tsi.mutex.RLock()
	defer tsi.mutex.RUnlock()

	copy(tweets, tsi.tweets[:])

	return tweets
}

func (tsi *tweetStreamInfo) storeTweets(conf *Config, sse sse.Provider) {
	for {
		c := <-sse.NewClients
		for m := range c.Channel {
			switch msg := m.(type) {
			default:
			case Tweet:
				tsi.Add(&msg)
			}
		}
	}
}

func (tsi tweetStreamInfo) ContainsTweet(t Tweet) bool {
	for _, hashtag := range tsi.Hashtags {
		for _, tweetHashtag := range t.Entities.Hashtags {
			if strings.EqualFold(hashtag, tweetHashtag.Text) {
				return true
			}
		}
	}
	for _, user := range tsi.Users {
		if strings.EqualFold(user, t.User.ScreenName) {
			return true
		}
	}
	return false
}
