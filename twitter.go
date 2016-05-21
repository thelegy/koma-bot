package main

import (
	"encoding/json"
	"net/url"
	"regexp"
	"time"

	"github.com/chimeracoder/anaconda"
)

func newTwitterApi(conf *Config) *anaconda.TwitterApi {

	anaconda.SetConsumerKey(conf.GetConfigString("twitter.login.consumer_key"))
	anaconda.SetConsumerSecret(conf.GetConfigString("twitter.login.consumer_secret"))

	api := anaconda.NewTwitterApi(conf.GetConfigString("twitter.login.access_token_key"),
		conf.GetConfigString("twitter.login.access_token_secret"))

	return api
}

func newTwitterStream(conf *Config, api *anaconda.TwitterApi) *anaconda.Stream {
	params := url.Values{}
	params.Set("track", conf.GetConfigString("twitter.track"))

	stream := api.PublicStreamFilter(params)

	return stream

}

func convertTweet(t anaconda.Tweet) (Tweet, error) {
	var tweet Tweet

	jsonTweet, err := json.Marshal(t)
	if err != nil {
		return Tweet{}, err
	}
	err = json.Unmarshal(jsonTweet, &tweet)
	if err != nil {
		return Tweet{}, err
	}

	return tweet, nil
}

func processStream(api *anaconda.TwitterApi, stream *anaconda.Stream, sseEventStream chan<- interface{}) {
	regex, err := regexp.Compile("https?://")
	for message := range stream.C {
		if t, ok := message.(anaconda.Tweet); ok {
			if err == nil && regex.MatchString(t.Text) {
				// Tweet might contain an image
				// we need to hydrate the tweet first
				go hydrateTweet(api, sseEventStream, t.Id)
				continue
			}
			tweet, err := convertTweet(t)
			if err != nil {
				// log here
				continue
			}
			sseEventStream <- tweet
		}
	}
}

func hydrateTweet(api *anaconda.TwitterApi, sseEventStream chan<- interface{}, tweetId int64) {
	t, err := api.GetTweet(tweetId, nil)
	if err != nil {
		// log here
		return
	}
	tweet, err := convertTweet(t)
	if err != nil {
		// log here
		return
	}
	sseEventStream <- tweet
}

func twitterListen(conf *Config, sseEventStream chan<- interface{}) *anaconda.TwitterApi {
	api := newTwitterApi(conf)
	if conf.IsDebugging() {
		api.SetLogger(anaconda.BasicLogger)
	}

	go func(api *anaconda.TwitterApi, conf *Config, sseEventStream chan<- interface{}) {
		for {
			stream := newTwitterStream(conf, api)
			processStream(api, stream, sseEventStream)

			//stream closed, need to wait & restart it
			<-time.After(60 * time.Second)
		}
	}(api, conf, sseEventStream)
	return api
}

func loadRecentTweets(api *anaconda.TwitterApi, conf *Config, ts *TweetStorage) error {
	val := url.Values{
		"count":            []string{"100"},
		"include_entities": []string{"true"},
	}
	tweets, err := api.GetSearch(conf.GetConfigString("twitter.track"), val)
	if err != nil {
		return err
	}

	for _, t := range tweets.Statuses {
		tweet, err := convertTweet(t)
		if err != nil {
			// log here
			continue
		}
		ts.Add(&tweet)
	}

	return nil
}
