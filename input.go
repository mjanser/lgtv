package lgtv

const (
	uriInputList = "ssap://tv/getExternalInputList"
	uriInputSet  = "ssap://tv/switchInput"
)

// Input contains information about a device input retrieved from the LG TV
type Input struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	Port         int    `json:"port"`
	AppID        string `json:"appId"`
	Icon         string `json:"icon"`
	Modified     bool   `json:"modified"`
	LastUniqueID int    `json:"lastUniqueId"`
	SubList      []struct {
		ID           string `json:"id"`
		UniqueID     int    `json:"uniqueId"`
		CECPDeviceID int    `json:"cecpDevId"`
		CECPNewType  int    `json:"cecpNewType"`
		Version      int    `json:"version"`
		OSDName      string `json:"osdName"`
	} `json:"subList"`
	SubCount  int  `json:"subCount"`
	Connected bool `json:"connected"`
	Favorite  bool `json:"favorite"`
}

// GetInputs returns a list of supported inputs
func (lgtv *Client) GetInputs() ([]Input, error) {
	resp := struct {
		Inputs []Input `json:"devices"`
	}{}
	if err := lgtv.c.Request(lgtv.nextID(), uriInputList, nil, &resp); err != nil {
		return resp.Inputs, err
	}

	return resp.Inputs, nil
}

// SwitchInput changes the active input
func (lgtv *Client) SwitchInput(inputID string) error {
	payload := struct {
		InputID string `json:"inputId"`
	}{
		InputID: inputID,
	}

	return lgtv.c.Request(lgtv.nextID(), uriInputSet, payload, nil)
}
