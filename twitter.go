package main

import (
	"encoding/json"
	"github.com/chimeracoder/anaconda"
	"net/url"
	"time"
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

func processStream(stream *anaconda.Stream, sseEventStream chan<- interface{}) {
	for message := range stream.C {
		if t, ok := message.(anaconda.Tweet); ok {
			tweet, err := convertTweet(t)
			if err != nil {
				// log here
				continue
			}
			sseEventStream <- tweet
		}
	}
}

func twitterListen(conf *Config, sseEventStream chan<- interface{}) {
	api := newTwitterApi(conf)
	if conf.IsDebugging() {
		api.SetLogger(anaconda.BasicLogger)
	}

	for {
		stream := newTwitterStream(conf, api)
		processStream(stream, sseEventStream)

		//stream closed, need to wait & restart it
		<-time.After(60 * time.Second)
	}
}
