package kahoot

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var ErrConnClosed = errors.New("connection closed")
var ErrNotSubscribed = errors.New("not subscribed to channel")

const incomingBufferSize = 16

type Message map[string]interface{}

type Conn struct {
	ws *websocket.Conn

	clientId string
	gameId   int

	channelsLock sync.RWMutex
	incoming     map[string]chan Message
	outgoing     chan Message

	closed chan struct{}
}

// NewConn connects to the kahoot server and performs a handshake
// using a given game pin.
func NewConn(gameId int) (*Conn, error) {
	conn, err := net.Dial("tcp", "kahoot.it:443")
	if err != nil {
		return nil, err
	}

	url, _ := url.Parse("wss://kahoot.it/cometd")
	reqHeader := http.Header{}
	reqHeader.Set("Origin", "https://kahoot.it")
	reqHeader.Set("Cookie", "no.mobitroll.session="+strconv.Itoa(gameId))
	ws, _, err := websocket.NewClient(conn, url, reqHeader, 100, 100)
	if err != nil {
		return nil, err
	}

	c := &Conn{
		ws:     ws,
		gameId: gameId,
		incoming: map[string]chan Message{
			"/meta/connect":    make(chan Message, incomingBufferSize),
			"/meta/disconnect": make(chan Message, incomingBufferSize),
			"/meta/handshake":  make(chan Message, incomingBufferSize),
			"/meta/subscribe":  make(chan Message, incomingBufferSize),
		},
		outgoing: make(chan Message),
		closed:   make(chan struct{}),
	}

	go c.readLoop()
	go c.writeLoop()

	err = c.Send("/meta/handshake", Message{
		"version":                  "1.0",
		"minimumVersion":           "1.0",
		"supportedConnectionTypes": []string{"websocket", "long-polling"},
		"advice":                   map[string]int{"timeout": 60000, "interval": 0},
	})
	if err != nil {
		c.Close()
		return nil, err
	}

	response, err := c.Receive("/meta/handshake")
	if err != nil {
		c.Close()
		return nil, err
	}

	if clientId, ok := response["clientId"].(string); !ok {
		c.Close()
		return nil, errors.New("invalid handshake response")
	} else {
		c.clientId = clientId
	}

	for _, service := range []string{"controller", "player", "status"} {
		if err := c.Subscribe("/service/" + service); err != nil {
			c.Close()
			return nil, err
		}
	}

	err = c.Send("/meta/connect", Message{
		"connectionType": "websocket",
		"advice":         map[string]int{"timeout": 0},
	})
	if err != nil {
		c.Close()
		return nil, err
	}
	connResp, err := c.Receive("/meta/connect")
	if err != nil {
		c.Close()
		return nil, err
	}
	if success, ok := connResp["successful"].(bool); !ok || !success {
		c.Close()
		return nil, errors.New("did not receive successful response")
	}

	go c.keepAliveLoop()

	return c, nil
}

// Login tells the server our nickname.
func (c *Conn) Login(nickname string) error {
	m := Message{
		"data": Message{
			"type":   "login",
			"gameid": c.gameId,
			"host":   "kahoot.it",
			"name":   nickname,
		},
	}
	if err := c.Send("/service/controller", m); err != nil {
		return err
	}

	for {
		resp, err := c.Receive("/service/controller")
		if err != nil {
			return err
		} else if data, ok := resp["data"].(map[string]interface{}); !ok {
			continue
		} else if typeStr, ok := data["type"].(string); !ok || typeStr != "loginResponse" {
			continue
		} else {
			return nil
		}
	}
}

// Close terminates the connection, waiting synchronously for the
// incoming channels to close.
func (c *Conn) Close() {
	c.ws.Close()
	<-c.closed
}

// GracefulClose closes the connection gracefully, telling the other end that
// we are disconnecting.
func (c *Conn) GracefulClose() {
	defer c.Close()
	if c.Send("/meta/disconnect", Message{}) != nil {
		return
	}
	c.Receive("/meta/disconnect")
}

// Send transmits a message to the server over a channel.
func (c *Conn) Send(channel string, m Message) error {
	packet := Message{}
	for k, v := range m {
		packet[k] = v
	}
	packet["channel"] = channel

	select {
	case c.outgoing <- packet:
	case <-c.closed:
		return ErrConnClosed
	}
	return nil
}

// Subscribe tells the server that we wish to receive messages
// on a given channel.
func (c *Conn) Subscribe(name string) error {
	c.channelsLock.Lock()
	if c.incoming == nil {
		c.channelsLock.Unlock()
		return ErrConnClosed
	} else if _, ok := c.incoming[name]; ok {
		c.channelsLock.Unlock()
		return nil
	}
	c.incoming[name] = make(chan Message, incomingBufferSize)
	c.channelsLock.Unlock()
	c.Send("/meta/subscribe", Message{"subscription": name})
	nextMsg, err := c.Receive("/meta/subscribe")
	if err != nil {
		return err
	} else if success, ok := nextMsg["successful"].(bool); !ok || !success {
		return errors.New("did not receive successful response")
	}
	return nil
}

// Receive returns the next message on a given channel.
// You must Subscribe() to the channel before Receiving on it.
func (c *Conn) Receive(channel string) (Message, error) {
	c.channelsLock.RLock()
	if c.incoming == nil {
		c.channelsLock.RUnlock()
		return nil, ErrConnClosed
	}
	ch, ok := c.incoming[channel]
	c.channelsLock.RUnlock()
	if !ok {
		return nil, ErrNotSubscribed
	}
	if res := <-ch; res != nil {
		return res, nil
	} else {
		return nil, ErrConnClosed
	}
}

func (c *Conn) readLoop() {
	defer func() {
		c.channelsLock.Lock()
		defer c.channelsLock.Unlock()
		for _, ch := range c.incoming {
			close(ch)
		}
		c.incoming = nil
		close(c.closed)
	}()
	for {
		var msgs []Message
		err := c.ws.ReadJSON(&msgs)
		if err != nil {
			return
		}
		for _, msg := range msgs {
			if chName, ok := msg["channel"].(string); !ok {
				return
			} else {
				c.channelsLock.RLock()
				ch, ok := c.incoming[chName]
				c.channelsLock.RUnlock()
				if ok {
					// NOTE: the select allows us to drop packets from channels
					// that nobody cares about (e.g. /meta/connect).
					select {
					case ch <- msg:
					default:
					}
				}
			}
		}
	}
}

func (c *Conn) writeLoop() {
	id := 0
	for {
		select {
		case msg := <-c.outgoing:
			id++
			msg["id"] = strconv.Itoa(id)
			if msg["channel"] != "/meta/handshake" {
				msg["clientId"] = c.clientId
			}
			if c.ws.WriteJSON([]Message{msg}) != nil {
				c.ws.Close()
				return
			}
		case <-c.closed:
			return
		}
	}
}

func (c *Conn) keepAliveLoop() {
	for {
		delay := time.After(time.Second * 5)
		select {
		case <-delay:
		case <-c.closed:
			return
		}
		c.Send("/meta/connect", Message{"connectionType": "websocket"})
	}
}
