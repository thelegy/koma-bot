package main

import (
	"encoding/json"
	"github.com/chimeracoder/anaconda"
	"net/url"
	"time"
)

func newTwitterApi() *anaconda.TwitterApi {

	anaconda.SetConsumerKey(getConfigString("twitter.login.consumer_key"))
	anaconda.SetConsumerSecret(getConfigString("twitter.login.consumer_secret"))

	api := anaconda.NewTwitterApi(getConfigString("twitter.login.access_token_key"),
		getConfigString("twitter.login.access_token_secret"))

	return api
}

func newTwitterStream(api *anaconda.TwitterApi) *anaconda.Stream {
	params := url.Values{}
	params.Set("track", getConfigString("twitter.track"))

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

func processStream(stream *anaconda.Stream) {
	for message := range stream.C {
		if t, ok := message.(anaconda.Tweet); ok {
			tweet, err := convertTweet(t)
			if err != nil {
				// log here
				continue
			}
			SSE.EventStream <- tweet
		}
	}
}

func twitterListen() {
	api := newTwitterApi()
	if Debug {
		api.SetLogger(anaconda.BasicLogger)
	}

	for {
		stream := newTwitterStream(api)
		processStream(stream)

		//stream closed, need to wait & restart it
		<-time.After(60 * time.Second)
	}
}
