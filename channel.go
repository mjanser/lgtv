package lgtv

import (
	"encoding/json"
	"log"
)

const (
	uriChannelGet  = "ssap://tv/getCurrentChannel"
	uriChannelSet  = "ssap://tv/openChannel"
	uriChannelList = "ssap://tv/getChannelList"
	uriChannelInfo = "ssap://tv/getChannelProgramInfo"
)

// Channel contains information about a channel retrieved from the LG TV
type Channel struct {
	ID              string `json:"channelID"`
	SignalChannelID string `json:"signalChannelId"`
	ModeID          int    `json:"channelModeId"`
	ModeName        string `json:"channelModeName"`
	TypeID          int    `json:"channelTypeId"`
	TypeName        string `json:"channelTypeName"`
	Number          string `json:"channelNumber"`
	Name            string `json:"channelName"`
	PhysicalNumber  int    `json:"physicalNumber"`
	Changed         bool   `json:"isChannelChanged"`
	Skipped         bool   `json:"isSkipped"`
	Locked          bool   `json:"isLocked"`
	Descrambled     bool   `json:"isDescrambled"`
	Scrambled       bool   `json:"isScrambled"`
	FineTuned       bool   `json:"isFineTuned"`
	Invisible       bool   `json:"isInvisible"`
	FavoriteGroup   string `json:"favoriteGroup"`
	HevcChannel     bool   `json:"isHEVCChannel"`
	HybridTvType    string `json:"hybridtvType"`
	Dualchannel     struct {
		ID       string `json:"dualChannelID"`
		TypeID   int    `json:"dualChannelTypeId"`
		TypeName string `json:"dualChannelTypeName"`
		Number   string `json:"dualChannelNumber"`
	} `json:"dualChannel"`
}

// GetCurrentChannel returns the currently active channel
func (lgtv *Client) GetCurrentChannel() (Channel, error) {
	resp := Channel{}
	if err := lgtv.c.Request(lgtv.nextID(), uriChannelGet, nil, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

// GetChannels returns a list of available channels
func (lgtv *Client) GetChannels() ([]Channel, error) {
	resp := struct {
		Channels []Channel `json:"channels"`
	}{}
	if err := lgtv.c.Request(lgtv.nextID(), uriChannelList, nil, &resp); err != nil {
		return resp.Channels, err
	}

	return resp.Channels, nil
}

// SetChannel opens the given channel
func (lgtv *Client) SetChannel(number string) error {
	payload := struct {
		Number string `json:"channelNumber"`
	}{
		Number: number,
	}

	return lgtv.c.Request(lgtv.nextID(), uriChannelSet, payload, nil)
}

// SubscribeChannel registers a subscription for channel updates
func (lgtv *Client) SubscribeChannel(h func(Channel)) error {
	return lgtv.subscribe(uriChannelGet, func(payload json.RawMessage) {
		p := Channel{}
		if err := json.Unmarshal(payload, &p); err != nil {
			log.Printf("Could not unmarshal channel update \"%s\": %s", string(payload), err)
		} else {
			h(p)
		}
	})
}
