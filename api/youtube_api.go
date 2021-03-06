package api

import (
	"fmt"
	"github.com/lintmx/dd-recorder/utils"
	"github.com/tidwall/gjson"
	"regexp"
)

/**
Thanks:
- https://github.com/ytdl-org/youtube-dl/blob/master/youtube_dl/extractor/youtube.py
- https://github.com/streamlink/streamlink/blob/master/src/streamlink/plugins/youtube.py
*/

var (
	youtubeLiveURL = "https://www.youtube.com/channel/%s/live"
)

// YouTubeLive youtube live api
type YouTubeLive struct {
	BaseAPI
	videoID string
}

// NewYouTubeLive return a youtubeLive struct
func NewYouTubeLive(base *BaseAPI) *YouTubeLive {
	youtubeLive := YouTubeLive{
		BaseAPI: *base,
	}
	regexURL := regexp.MustCompile(`^(?:https?:\/\/)?www\.youtube\.com\/channel\/([^\/]+)(?:[\/])?(?:live)?`)
	if result := regexURL.FindStringSubmatch(youtubeLive.GetLiveURL()); result != nil {
		youtubeLive.liveID = result[1]
		return &youtubeLive
	}

	return nil
}

// get youtube live page js
func (y *YouTubeLive) getLiveInfo() (string, error) {
	body, err := utils.HTTPGet(fmt.Sprintf(youtubeLiveURL, y.liveID))

	if err != nil {
		return "", fmt.Errorf("Http Error - youtubeLiveURL - %s", err.Error())
	} else if body == "" {
		return "", fmt.Errorf("youtubeLiveURL download failed")
	}

	liveStatusRegex := []*regexp.Regexp{
		regexp.MustCompile(`;ytplayer\.config\s*=\s*({.+?});ytplayer`),
		regexp.MustCompile(`;ytplayer\.config\s*=\s*({.+?});`),
	}

	for _, re := range liveStatusRegex {
		data := re.FindStringSubmatch(body)

		if len(data) >= 2 {
			return data[1], nil
		}
	}

	// TODO: get live info by video id

	return "", fmt.Errorf("youtubeLive get live info error")
}

// RefreshLiveInfo refresh live info
func (y *YouTubeLive) RefreshLiveInfo() error {
	liveInfo, err := y.getLiveInfo()

	if err != nil {
		return err
	}

	status := gjson.Get(liveInfo, "args.livestream")

	if !status.Exists() || status.Int() != 1 {
		y.liveStatus = false
		return nil
	}

	y.liveStatus = true
	y.liveAuthor = gjson.Get(liveInfo, "args.author").String()
	y.liveTitle = gjson.Get(liveInfo, "args.title").String()

	return nil
}
