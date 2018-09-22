package lgtv

import (
	"encoding/json"
	"log"
)

const (
	uriVolumeGet = "ssap://audio/getVolume"
	uriVolumeSet = "ssap://audio/setVolume"
	uriMuteGet   = "ssap://audio/getMute"
	uriMuteSet   = "ssap://audio/setMute"
)

// Volume contains information about the volume retrieved from the LG TV
type Volume struct {
	Muted     bool   `json:"muted"`
	Scenario  string `json:"scenario,omitempty"`
	Active    bool   `json:"active,omitempty"`
	Action    string `json:"action,omitempty"`
	Volume    int    `json:"volume"`
	MaxVolume int    `json:"volumeMax,omitempty"`
}

// GetVolume returns the current volume information
func (lgtv *Client) GetVolume() (Volume, error) {
	resp := Volume{}
	if err := lgtv.c.Request(lgtv.nextID(), uriVolumeGet, nil, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

// SetVolume sets the volume of the LG TV
func (lgtv *Client) SetVolume(v int) error {
	payload := struct {
		Volume int `json:"volume"`
	}{
		Volume: v,
	}
	return lgtv.c.Request(lgtv.nextID(), uriVolumeSet, payload, nil)
}

// IsMute returns the mute state
func (lgtv *Client) IsMute() (bool, error) {
	resp := struct {
		Mute bool `json:"mute"`
	}{}
	if err := lgtv.c.Request(lgtv.nextID(), uriMuteGet, nil, &resp); err != nil {
		return false, err
	}

	return resp.Mute, nil
}

// SetMute sets the mute state of the LG TV
func (lgtv *Client) SetMute(m bool) error {
	payload := struct {
		Mute bool `json:"mute"`
	}{
		Mute: m,
	}
	return lgtv.c.Request(lgtv.nextID(), uriMuteSet, payload, nil)
}

// SubscribeVolume registers a subscription for volume updates
func (lgtv *Client) SubscribeVolume(h func(Volume)) error {
	return lgtv.subscribe(uriVolumeGet, func(payload json.RawMessage) {
		p := Volume{}
		if err := json.Unmarshal(payload, &p); err != nil {
			log.Printf("Could not unmarshal volume update \"%s\": %s", string(payload), err)
		} else {
			h(p)
		}
	})
}
