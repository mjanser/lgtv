package lgtv

import (
	"encoding/json"
	"sync"

	"github.com/mjanser/lgtv/connection"
)

var permissions = []string{
	"LAUNCH",
	"CONTROL_AUDIO",
	"CONTROL_POWER",
	"CONTROL_INPUT_TV",
	"CONTROL_INPUT_MEDIA_PLAYBACK",
	"READ_TV_CHANNEL_LIST",
	"READ_CURRENT_CHANNEL",
	"READ_RUNNING_APPS",
	"READ_INSTALLED_APPS",
	"READ_INPUT_DEVICE_LIST",
	"WRITE_NOTIFICATION_TOAST",
}

// Client represents a LG TV client
type Client struct {
	c             *connection.Connection
	key           string
	responseChans map[int]chan []byte
	requestIDLock sync.Mutex
	requestID     int
	onConnect     func(*Client)
	onDisconnect  func(*Client, error)
}

// NewDefaultClient creates a LG TV client with default timeouts
func NewDefaultClient(url string, key string) *Client {
	return NewClient(url, key, connection.DefaultTimeouts())
}

// NewClient creates a LG TV client
func NewClient(url string, key string, t connection.Timeouts) *Client {
	return &Client{
		c:             connection.NewConnection(url, t),
		key:           key,
		responseChans: make(map[int]chan []byte),
	}
}

// OnConnect registers a callback which is executed when a connection was established
func (lgtv *Client) OnConnect(f func(*Client)) {
	lgtv.onConnect = f
}

// OnDisconnect registers a callback which is executed when the connection was closed
func (lgtv *Client) OnDisconnect(f func(*Client, error)) {
	lgtv.onDisconnect = f
	lgtv.c.OnDisconnect = func(c *connection.Connection, err error) {
		lgtv.Disconnect()
	}
}

// Connect tries to connect to the LG TV
func (lgtv *Client) Connect() error {
	if lgtv.IsConnected() {
		return nil
	}
	if err := lgtv.c.Connect(); err != nil {
		return err
	}

	if err := lgtv.register(); err != nil {
		lgtv.Disconnect()

		return err
	}

	if lgtv.onConnect != nil {
		lgtv.onConnect(lgtv)
	}

	return nil
}

// Disconnect closes the connection to the LG TV
func (lgtv *Client) Disconnect() error {
	if !lgtv.IsConnected() {
		return nil
	}
	if err := lgtv.c.Close(); err != nil {
		return err
	}

	if lgtv.onDisconnect != nil {
		lgtv.onDisconnect(lgtv, nil)
	}

	return nil
}

// IsConnected returns whether or not there is an open connection to the LG TV
func (lgtv *Client) IsConnected() bool {
	return lgtv.c.IsConnected()
}

func (lgtv *Client) register() error {
	return lgtv.c.Register(lgtv.nextID(), lgtv.key, permissions)
}

func (lgtv *Client) subscribe(uri string, callback func(json.RawMessage)) error {
	ch, err := lgtv.c.Subscribe(lgtv.nextID(), uri)
	if err != nil {
		return err
	}

	go func() {
		for lgtv.IsConnected() {
			select {
			case resp := <-ch:
				callback(resp.Payload)
			}
		}
	}()

	return nil
}

func (lgtv *Client) nextID() int {
	lgtv.requestIDLock.Lock()
	defer lgtv.requestIDLock.Unlock()

	lgtv.requestID++
	return lgtv.requestID
}
