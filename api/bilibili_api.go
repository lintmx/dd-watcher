package api

import (
	"fmt"
	"github.com/lintmx/dd-recorder/utils"
	"github.com/tidwall/gjson"
	"regexp"
)

var (
	bilibiliRealRoomIDAPI = "https://api.live.bilibili.com/room/v1/Room/room_init?id=%s"
	bilibiliRoomInfoAPI   = "https://api.live.bilibili.com/room/v1/Room/get_info?room_id=%s"
	bilibiliRoomAnchorAPI = "https://api.live.bilibili.com/live_user/v1/UserInfo/get_anchor_in_room?roomid=%s"
)

// BilibiliLive bilibili live api
type BilibiliLive struct {
	BaseAPI
	roomID string
}

// NewBilibiliLive return a bilibililive struct
func NewBilibiliLive(base *BaseAPI) *BilibiliLive {
	bilibiliLive := BilibiliLive{
		BaseAPI: *base,
	}
	regexURL := regexp.MustCompile(`^(?:https?:\/\/)?live\.bilibili\.com\/(\d+)[\/\?\#]?.*$`)
	if result := regexURL.FindStringSubmatch(bilibiliLive.GetLiveURL()); result != nil {
		bilibiliLive.liveID = result[1]
		return &bilibiliLive
	}

	return nil
}

func (b *BilibiliLive) getRealRoomID() error {
	body, err := utils.HTTPGet(fmt.Sprintf(bilibiliRealRoomIDAPI, b.liveID))

	if err != nil {
		return fmt.Errorf("Http Error - bilibiliRealRoomIDAPI - %s", err.Error())
	} else if code := gjson.Get(body, "code"); !code.Exists() {
		return fmt.Errorf("bilibiliRealRoomIDAPI is broken")
	} else if code.Int() != 0 {
		return fmt.Errorf("bilibiliRealRoomIDAPI - %s", gjson.Get(body, "msg").String())
	}

	b.roomID = gjson.Get(body, "data.room_id").String()

	return nil
}

// RefreshLiveInfo refresh live info
func (b *BilibiliLive) RefreshLiveInfo() error {
	if b.roomID == "" {
		if err := b.getRealRoomID(); err != nil {
			return err
		}
	}

	// get live title and live status
	body, err := utils.HTTPGet(fmt.Sprintf(bilibiliRoomInfoAPI, b.roomID))

	if err != nil {
		return fmt.Errorf("Http Error - bilibiliRoomInfoAPI - %s", err.Error())
	} else if code := gjson.Get(body, "code"); !code.Exists() {
		return fmt.Errorf("bilibiliRoomInfoAPI is broken")
	} else if code.Int() != 0 {
		return fmt.Errorf("bilibiliRoomInfoAPI - %s", gjson.Get(body, "msg").String())
	}

	status := gjson.Get(body, "data.live_status").Int() == 1
	// return when live status not change
	if b.liveStatus == status {
		return nil
	}

	b.liveStatus = status
	b.liveTitle = gjson.Get(body, "data.title").String()

	// get live author
	body, err = utils.HTTPGet(fmt.Sprintf(bilibiliRoomAnchorAPI, b.roomID))

	if err != nil {
		return fmt.Errorf("Http Error - bilibiliRoomAnchorAPI - %s", err.Error())
	} else if code := gjson.Get(body, "code"); !code.Exists() {
		return fmt.Errorf("bilibiliRoomAnchorAPI is broken")
	} else if code.Int() != 0 {
		return fmt.Errorf("bilibiliRoomAnchorAPI - %s", gjson.Get(body, "msg").String())
	}

	b.liveAuthor = gjson.Get(body, "data.info.uname").String()

	return nil
}
