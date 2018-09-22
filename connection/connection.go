package connection

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Connection represents a websocket connection to a LG TV
type Connection struct {
	url               string
	ws                *websocket.Conn
	timeouts          Timeouts
	responseChans     map[int]chan Response
	notificationChans map[int]chan Response
	done              chan bool
	mutex             *sync.Mutex
	OnConnect         func(*Connection)
	OnDisconnect      func(*Connection, error)
}

// NewConnection creates a new connection to a LG TV with the specified timeouts
func NewConnection(url string, t Timeouts) *Connection {
	return &Connection{
		url:               url,
		timeouts:          t,
		responseChans:     make(map[int]chan Response),
		notificationChans: make(map[int]chan Response),
		mutex:             &sync.Mutex{},
	}
}

// Connect tries to establish a connection to the LG TV
func (c *Connection) Connect() error {
	if c.ws != nil {
		return nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	ws, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return err
	}

	c.ws = ws
	if err := c.ws.SetReadDeadline(time.Now().Add(c.timeouts.Read)); err != nil {
		c.disconnect(err)
		return err
	}
	c.ws.SetPongHandler(c.pong)

	c.done = make(chan bool)

	go c.listen()
	go c.ping()

	if c.OnConnect != nil {
		c.OnConnect(c)
	}

	return nil
}

// Close sends a close control message and disconnects from the LG TV
func (c *Connection) Close() error {
	if c.ws == nil {
		return errors.New("Cannot close connection because there is no connection")
	}

	if err := c.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		return err
	}

	c.disconnect(nil)
	return nil
}

// IsConnected returns whether or not a connection to the LG TV is established
func (c *Connection) IsConnected() bool {
	return c.ws != nil
}

func (c *Connection) disconnect(err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	log.Println("disconnect....")
	if c.ws == nil {
		log.Println("disconnect!")
		return
	}

	close(c.done)
	c.ws.Close()

	log.Println("disconnected....")

	c.ws = nil
	if c.OnDisconnect != nil {
		c.OnDisconnect(c, err)
	}
}
