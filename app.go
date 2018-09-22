package lgtv

import (
	"encoding/json"
	"log"
)

const (
	// AppIDLiveTV defines the app ID for Live TV
	AppIDLiveTV = "com.webos.app.livetv"
	// AppIDHDMI1 defines the app ID for HDMI 1 input
	AppIDHDMI1 = "com.webos.app.hdmi1"
	// AppIDHDMI2 defines the app ID for HDMI 2 input
	AppIDHDMI2 = "com.webos.app.hdmi2"
	// AppIDHDMI3 defines the app ID for HDMI 3 input
	AppIDHDMI3 = "com.webos.app.hdmi3"
	// AppIDHDMI4 defines the app ID for HDMI 4 input
	AppIDHDMI4 = "com.webos.app.hdmi4"

	uriAppGet    = "ssap://com.webos.applicationManager/getForegroundAppInfo"
	uriAppList   = "ssap://com.webos.applicationManager/listApps"
	uriAppLaunch = "ssap://system.launcher/launch"
)

// App contains information about an app retrieved from the LG TV
type App struct {
	ID                         string          `json:"id"`
	Title                      string          `json:"title"`
	Icon                       string          `json:"icon"`
	Visible                    bool            `json:"visible"`
	Type                       string          `json:"type"`
	DefaultWindowType          string          `json:"defaultWindowType"`
	InstalledTime              int             `json:"installedTime"`
	BgImages                   []string        `json:"bgImages"`
	UIRevision                 json.RawMessage `json:"uiRevision"`
	CPApp                      bool            `json:"CPApp"`
	Version                    string          `json:"version"`
	SystemApp                  bool            `json:"systemApp"`
	AppSize                    int             `json:"appsize"`
	Vendor                     string          `json:"vendor"`
	MiniIcon                   string          `json:"miniicon"`
	HasPromotion               bool            `json:"hasPromotion"`
	TileSize                   string          `json:"tileSize"`
	Icons                      []string        `json:"icons"`
	RequestedWindowOrientation string          `json:"requestedWindowOrientation"`
	LargeIcon                  string          `json:"largeIcon"`
	Lockable                   bool            `json:"lockable"`
	Transparent                bool            `json:"transparent"`
	CheckUpdateOnLaunch        bool            `json:"checkUpdateOnLaunch"`
	Category                   string          `json:"category"`
	LaunchInNewGroup           bool            `json:"launchinnewgroup"`
	SpinnerOnLaunch            bool            `json:"spinnerOnLaunch"`
	HandlesRelaunch            bool            `json:"handlesRelaunch"`
	UnMovable                  bool            `json:"unmovable"`
	Inspectable                bool            `json:"inspectable"`
	InAppSetting               bool            `json:"inAppSetting"`
	PriviledgedJail            bool            `json:"privilegedJail"`
	SupportQuickStart          bool            `json:"supportQuickStart"`
	SplashBackground           string          `json:"splashBackground"`
	TrustLevel                 string          `json:"trustLevel"`
	BootLaunchParams           struct {
		Boot bool `json:"boot"`
	} `json:"bootLaunchParams"`
	HardwareFeaturesNeeded int  `json:"hardwareFeaturesNeeded"`
	NoWindow               bool `json:"noWindow"`
	Age                    int  `json:"age"`
	WindowGroup            struct {
		Owner     bool `json:"owner"`
		OwnerInfo struct {
			AllowAnonymous bool `json:"allowAnonymous"`
			Layers         []struct {
				Z    int    `json:"z"`
				Name string `json:"name"`
			} `json:"layers"`
		} `json:"ownerInfo"`
		Name string `json:"name"`
	} `json:"windowGroup"`
	Accessibility struct {
		SupportsAudioGuidance bool `json:"supportsAudioGuidance"`
	} `json:"accessibility"`
	FolderPath            string `json:"folderPath"`
	DeepLinkingParams     string `json:"deeplinkingParams"`
	Main                  string `json:"main"`
	Removable             bool   `json:"removable"`
	BgImage               string `json:"bgImage"`
	IconColor             string `json:"iconColor"`
	DisableBackHistoryAPI bool   `json:"disableBackHistoryAPI"`
	NoSplashOnLaunch      bool   `json:"noSplashOnLaunch"`
}

// AppInfo represents basic information about an app retrieved on subscription
type AppInfo struct {
	ID        string `json:"appId"`
	WindowID  string `json:"windowId"`
	ProcessID string `json:"processId"`
}

// GetApps returns a list of installed apps
func (lgtv *Client) GetApps() ([]App, error) {
	resp := struct {
		Apps []App `json:"apps"`
	}{}
	if err := lgtv.c.Request(lgtv.nextID(), uriAppList, nil, &resp); err != nil {
		return resp.Apps, err
	}

	return resp.Apps, nil
}

// LaunchApp start the given app
func (lgtv *Client) LaunchApp(id string) error {
	payload := struct {
		ID         string      `json:"id"`
		ContentID  string      `json:"contentId"`
		Parameters interface{} `json:"params"`
	}{
		ID: id,
	}

	return lgtv.c.Request(lgtv.nextID(), uriAppLaunch, payload, nil)
}

// SubscribeApp registers a subscription for app updates
func (lgtv *Client) SubscribeApp(h func(AppInfo)) error {
	return lgtv.subscribe(uriAppGet, func(payload json.RawMessage) {
		p := AppInfo{}
		if err := json.Unmarshal(payload, &p); err != nil {
			log.Printf("Could not unmarshal app info update \"%s\": %s", string(payload), err)
		} else {
			h(p)
		}
	})
}
