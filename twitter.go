package main

import (
	"encoding/json"
	"net/url"
	"regexp"
	"time"

	"github.com/chimeracoder/anaconda"
)

type tweetStream struct {
	Hashtags []string
	Users    []string
}

func newTwitterApi(conf *Config) *anaconda.TwitterApi {

	anaconda.SetConsumerKey(conf.GetConfigString("twitter.login.consumer_key"))
	anaconda.SetConsumerSecret(conf.GetConfigString("twitter.login.consumer_secret"))

	api := anaconda.NewTwitterApi(conf.GetConfigString("twitter.login.access_token_key"),
		conf.GetConfigString("twitter.login.access_token_secret"))

	return api
}

func newTwitterStream(conf *Config, api *anaconda.TwitterApi) *anaconda.Stream {
	track, follow := conf.GetTweetFilter()

	params := url.Values{}
	params.Set("track", track)
	params.Set("follow", follow)

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

func processStream(conf *Config, api *anaconda.TwitterApi, stream *anaconda.Stream, sseEventStream chan<- interface{}) {
	regex, rerr := regexp.Compile("https?://")
	for message := range stream.C {
		if t, ok := message.(anaconda.Tweet); ok {
			if rerr == nil && regex.MatchString(t.Text) {
				// Tweet might contain an image
				// we need to hydrate the tweet first
				go hydrateTweet(conf, api, sseEventStream, t.Id)
				continue
			}
			tweet, err := convertTweet(t)
			if err != nil {
				// log here
				continue
			}
			conf.StoreTweet(tweet)
			sseEventStream <- tweet
		}
	}
}

func hydrateTweet(conf *Config, api *anaconda.TwitterApi, sseEventStream chan<- interface{}, tweetId int64) {
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
	conf.StoreTweet(tweet)
	sseEventStream <- tweet
}

func twitterConnect(conf *Config) *anaconda.TwitterApi {
	api := newTwitterApi(conf)
	if conf.IsDebugging() {
		api.SetLogger(anaconda.BasicLogger)
	}
	return api
}

func twitterListen(api *anaconda.TwitterApi, conf *Config, sseEventStream chan<- interface{}) {
	for {
		stream := newTwitterStream(conf, api)
		processStream(conf, api, stream, sseEventStream)

		//stream closed, need to wait & restart it
		<-time.After(60 * time.Second)
	}
}

func loadRecentTweets(api *anaconda.TwitterApi, conf *Config) error {
	val := url.Values{
		"count":            []string{"100"},
		"include_entities": []string{"true"},
	}
	for _, tsi := range conf.StreamInfo {
		tweets, err := api.GetSearch(tsi.GetTweetFilter(), val)
		if err != nil {
			return err
		}

		for _, t := range tweets.Statuses {
			tweet, err := convertTweet(t)
			if err != nil {
				// log here
				continue
			}
			tsi.Add(&tweet)
		}
	}

	return nil
}
