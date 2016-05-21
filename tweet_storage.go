package main

import (
	"errors"
	"sync"

	"github.com/kerwindena/koma-bot/sse"
)

const TweetStorageMaxCapacity = 100

type TweetStorage struct {
	tweets [TweetStorageMaxCapacity]*Tweet
	mutex  *sync.RWMutex
}

func newTweetStorage() *TweetStorage {
	ts := &TweetStorage{
		mutex: &sync.RWMutex{},
	}
	return ts
}

func (ts *TweetStorage) Add(t *Tweet) error {
	j := -1

	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	for j < TweetStorageMaxCapacity-1 {
		if ts.tweets[j+1] == nil {
			j++
			continue
		}

		if t.Id == ts.tweets[j+1].Id {
			return errors.New("Tweet already exists")
		}

		if t.Id < ts.tweets[j+1].Id {
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
		ts.tweets[i] = ts.tweets[i+1]
	}
	ts.tweets[j] = t

	return nil
}

func (ts *TweetStorage) getTweets() []*Tweet {
	tweets := make([]*Tweet, TweetStorageMaxCapacity)

	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	copy(tweets, ts.tweets[:])

	return tweets
}

func (ts *TweetStorage) storeTweets(conf *Config, sse sse.Provider) {
	for {
		c := <-sse.NewClients
		for m := range c.Channel {
			switch msg := m.(type) {
			default:
			case Tweet:
				ts.Add(&msg)
			}
		}
	}
}
