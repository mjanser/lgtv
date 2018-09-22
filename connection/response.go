package connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	responseTypeError = "error"
)

// Response represents websocket response from the LG TV
type Response struct {
	ID      int             `json:"id"`
	Type    string          `json:"type"`
	Error   string          `json:"error"`
	Payload json.RawMessage `json:"payload"`
}

func (r *Response) isError() bool {
	return r.Type == responseTypeError
}

func (c *Connection) waitForResponse(id int, result interface{}, timeout time.Duration) (*Response, error) {
	if !c.IsConnected() {
		return nil, errors.New("Cannot send message because there is no connection")
	}

	ch, ok := c.responseChans[id]
	if !ok {
		return nil, fmt.Errorf("No response channel found for ID %d", id)
	}

	select {
	case <-time.After(timeout):
		return nil, fmt.Errorf("Request timed out after %s", timeout)
	case resp := <-ch:
		if resp.isError() {
			return &resp, fmt.Errorf("Request returned an error \"%s\"", resp.Error)
		}

		if result == nil {
			return &resp, nil
		}

		if err := json.Unmarshal(resp.Payload, result); err != nil {
			return &resp, fmt.Errorf("Could not unmarshal payload of response (%s): %s", string(resp.Payload), err)
		}

		return &resp, nil
	}
}

func (c *Connection) addResponseChannel(id int) chan Response {
	ch := make(chan Response)
	c.responseChans[id] = ch

	return ch
}

func (c *Connection) removeResponseChannel(id int) {
	if ch, ok := c.responseChans[id]; ok {
		close(ch)
		delete(c.responseChans, id)
	}
}
