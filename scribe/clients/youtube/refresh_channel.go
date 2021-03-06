package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type VideoSnippet struct {
	PublishedAt          time.Time         `json:"publishedAt"`
	ChannelID            string            `json:"channelId"`
	Title                string            `json:"title"`
	Description          string            `json:"description"`
	ChannelTitle         string            `json:"channelTitle"`
	LiveBroadcastContent string            `json:"liveBroadcastContent"`
	Thumbnails           SnippetThumbnails `json:"thumbnails"`
}

type VideoItemID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

type VideoItem struct {
	Kind    string       `json:"kind"`
	Etag    string       `json:"etag"`
	ID      VideoItemID  `json:"id"`
	Snippet VideoSnippet `json:"snippet"`
}

func (y *YoutubeClient) FetchVideos(channelID string) (FetchResponse[VideoItem], error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?channelId=%s&part=snippet&key=%s&maxResults=50&order=date", channelID, y.apiKey)

	response, err := http.Get(url) //#nosec G107 -- False positive
	if err != nil {
		return FetchResponse[VideoItem]{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return FetchResponse[VideoItem]{}, err
	}

	var data FetchResponse[VideoItem]
	if err := json.Unmarshal(responseData, &data); err != nil {
		return FetchResponse[VideoItem]{}, err
	}

	return data, nil
}
