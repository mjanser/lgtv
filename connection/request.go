package connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

const (
	requestTypeRequest   = "request"
	requestTypeSubscribe = "subscribe"
)

type request struct {
	ID      int         `json:"id"`
	Type    string      `json:"type"`
	URI     string      `json:"uri,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

// Request sends a request to the LG TV and waits for the response
func (c *Connection) Request(id int, uri string, payload interface{}, result interface{}) error {
	c.addResponseChannel(id)
	defer c.removeResponseChannel(id)

	if err := c.sendRequest(id, requestTypeRequest, uri, payload); err != nil {
		return err
	}

	_, err := c.waitForResponse(id, result, c.timeouts.Response)
	if err != nil {
		return fmt.Errorf("Request with ID %d and URI %s to LG TV at %s returned an error \"%s\"", id, uri, c.url, err)
	}

	return nil
}

// Subscribe sends a subscription request to the LG TV
func (c *Connection) Subscribe(id int, uri string) (chan Response, error) {
	ch := make(chan Response)
	c.notificationChans[id] = ch

	if err := c.sendRequest(id, requestTypeSubscribe, uri, nil); err != nil {
		close(ch)
		delete(c.notificationChans, id)
		return nil, err
	}

	return ch, nil
}

func (c *Connection) sendRequest(id int, reqType string, uri string, payload interface{}) error {
	if !c.IsConnected() {
		return errors.New("Cannot send message because there is no connection")
	}

	req := request{
		ID:      id,
		URI:     uri,
		Type:    reqType,
		Payload: payload,
	}

	if err := c.ws.WriteJSON(req); err != nil {
		return fmt.Errorf("Could not send request with ID %d, type %s and URI %s, error was \"%s\"", id, reqType, uri, err)
	}

	msg, _ := json.Marshal(req)
	log.Printf("-> %#v", string(msg))

	return nil
}
