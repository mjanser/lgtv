package lgtv

const (
	uriToast    = "ssap://system.notifications/createToast"
	uriPowerOff = "ssap://system/turnOff"
)

// Toast displays a message on the LG TV
func (lgtv *Client) Toast(message string) error {
	payload := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}

	return lgtv.c.Request(lgtv.nextID(), uriToast, payload, nil)
}

// TurnOff switches the TV off
func (lgtv *Client) TurnOff() error {
	return lgtv.c.Request(lgtv.nextID(), uriPowerOff, nil, nil)
}
