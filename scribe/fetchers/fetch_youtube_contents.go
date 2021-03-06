package fetchers

import (
	"fmt"

	"github.com/k3a/html2text"
	"github.com/skatekrak/scribe/clients/youtube"
)

func (fe *Fetcher) FetchYoutubeChannelContents(channelID string) ([]ContentFetchData, error) {
	data, err := fe.y.FetchVideos(channelID)
	if err != nil {
		return []ContentFetchData{}, err
	}

	items := make([]ContentFetchData, len(data.Items))

	for i, item := range data.Items {
		items[i] = ContentFetchData{
			Title:          html2text.HTML2Text(item.Snippet.Title),
			Description:    html2text.HTML2Text(item.Snippet.Description),
			PublishedAt:    item.Snippet.PublishedAt,
			RawDescription: item.Snippet.Description,
			ThumbnailURL:   youtube.GetBestThumbnail(item.Snippet.Thumbnails),
			ContentID:      item.ID.VideoID,
			ContentURL:     fmt.Sprintf("https://youtube.com/watch?=%s", item.ID.VideoID),
		}
	}

	return items, nil
}

func (fe *Fetcher) FetchYoutubeContent(channelIDs []string, contents map[string][]ContentFetchData) map[string]error {
	errors := make(map[string]error)

	for _, channelID := range channelIDs {
		c, err := fe.FetchYoutubeChannelContents(channelID)
		if err != nil {
			errors[channelID] = err
		} else {
			contents[channelID] = c
		}
	}

	return errors
}
