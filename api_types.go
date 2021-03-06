package main

import (
	"github.com/chimeracoder/anaconda"
)

const (
	MessageNone   = "none"
	MessageReload = "reload"
	MessageTweet  = "tweet"
	MessageSound  = "sound"
)

type TUser struct {
	CreatedAt string `json:"created_at"`
	// DefaultProfile                 bool     `json:"default_profile"`
	// DefaultProfileImage            bool     `json:"default_profile_image"`
	// Description                    string   `json:"description"`
	// FavouritesCount                int      `json:"favourites_count"`
	// FollowersCount                 int      `json:"followers_count"`
	Id int64 `json:"id"`
	// IdStr                          string   `json:"id_str"`
	// IsTranslator                   bool     `json:"is_translator"`
	// Lang                           string   `json:"lang"`
	Name string `json:"name"`
	// ProfileBackgroundColor         string   `json:"profile_background_color"`
	// ProfileBackgroundImageURL      string   `json:"profile_background_image_url"`
	// ProfileBackgroundImageUrlHttps string   `json:"profile_background_image_url_https"`
	// ProfileBackgroundTile          bool     `json:"profile_background_tile"`
	// ProfileBannerURL               string   `json:"profile_banner_url"`
	// ProfileImageURL                string   `json:"profile_image_url"`
	ProfileImageUrlHttps string `json:"profile_image_url_https"`
	// ProfileLinkColor               string   `json:"profile_link_color"`
	// ProfileSidebarBorderColor      string   `json:"profile_sidebar_border_color"`
	// ProfileSidebarFillColor        string   `json:"profile_sidebar_fill_color"`
	// ProfileTextColor               string   `json:"profile_text_color"`
	// ProfileUseBackgroundImage      bool     `json:"profile_use_background_image"`
	// Protected                      bool     `json:"protected"`
	ScreenName string `json:"screen_name"`
	// StatusesCount                  int64    `json:"statuses_count"`
	// TimeZone                       string   `json:"time_zone"`
	// URL                            string   `json:"url"`
	// UtcOffset                      int      `json:"utc_offset"`
	// Verified                       bool     `json:"verified"`
	// WithheldInCountries            []string `json:"withheld_in_countries"`
	// WithheldScope                  string   `json:"withheld_scope"`
}

type TEntityMedia struct {
	Id int64
	// Id_str               string
	Media_url       string
	Media_url_https string
	// Url                  string
	// Display_url          string
	// Expanded_url         string
	Sizes anaconda.MediaSizes
	// Source_status_id     int64
	// Source_status_id_str string
	Type    string
	Indices []int
	// VideoInfo            anaconda.VideoInfo `json:"video_info"`
}

type TEntities struct {
	Hashtags []struct {
		Indices []int
		Text    string
	}
	Urls []struct {
		Indices      []int
		Url          string
		Display_url  string
		Expanded_url string
	}
	Url           anaconda.UrlEntity
	User_mentions []struct {
		Name        string
		Indices     []int
		Screen_name string
		Id          int64
		Id_str      string
	}
	Media []TEntityMedia
}

type Tweet struct {
	CreatedAt     string `json:"created_at"`
	FavoriteCount int    `json:"favorite_count"`
	// FilterLevel          string   `json:"filter_level"`
	Id       int64     `json:"id"`
	Entities TEntities `json:"entities"`
	// ExtendedEntities     Entities `json:"extended_entities"`
	// IdStr                string   `json:"id_str"`
	// InReplyToScreenName  string   `json:"in_reply_to_screen_name"`
	// InReplyToStatusID    int64    `json:"in_reply_to_status_id"`
	// InReplyToStatusIdStr string   `json:"in_reply_to_status_id_str"`
	// InReplyToUserID      int64    `json:"in_reply_to_user_id"`
	// InReplyToUserIdStr   string   `json:"in_reply_to_user_id_str"`
	// Lang                 string   `json:"lang"`
	// RetweetCount         int      `json:"retweet_count"`
	// Retweeted            bool     `json:"retweeted"`
	RetweetedStatus *Tweet `json:"retweeted_status"`
	// Source               string   `json:"source"`
	// Text      string `json:"text"`
	FullText  string `json:"full_text"`
	Truncated bool   `json:"truncated"`
	User      TUser  `json:"user"`
	// WithheldCopyright    bool     `json:"withheld_copyright"`
	// WithheldInCountries  []string `json:"withheld_in_countries"`
	// WithheldScope        string   `json:"withheld_scope"`
}
