package kahoot

import (
	"crypto/tls"
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
)

type Connection struct {
	ws       *websocket.Conn
	lastId   int32
	gameId   int
	clientId string
	incoming chan incomingData
}

type Packet struct {
	Channel string
	Id      string
	Content map[string]interface{}
}

type PacketFilter func(*Packet) bool

type incomingData struct {
	packet *Packet
	err    error
}

func NewConnection(pin string) (*Connection, error) {
	gameId, err := strconv.Atoi(pin)
	if err != nil {
		return nil, err
	}

	conn, err := tls.Dial("tcp", "kahoot.it:443", nil)
	if err != nil {
		return nil, err
	}

	url, _ := url.Parse("wss://kahoot.it/cometd")
	reqHeader := http.Header{}
	reqHeader.Set("Origin", "https://kahoot.it")
	reqHeader.Set("Cookie", "no.mobitroll.session="+pin)
	ws, _, err := websocket.NewClient(conn, url, reqHeader, 100, 100)
	if err != nil {
		return nil, err
	}

	incoming := make(chan incomingData, 100)
	res := &Connection{ws, 0, gameId, "", incoming}
	go res.readLoop()
	return res, nil
}

func (c *Connection) Packet(channel string, d map[string]interface{}) *Packet {
	nextId := int(atomic.AddInt32(&c.lastId, 1))
	return &Packet{channel, strconv.Itoa(nextId), d}
}

func (c *Connection) Read() (*Packet, error) {
	if val, ok := <-c.incoming; !ok {
		return nil, errors.New("Connection closed.")
	} else {
		return val.packet, val.err
	}
}

func (c *Connection) Write(p *Packet) error {
	obj := map[string]interface{}{"id": p.Id, "channel": p.Channel,
		"ext": map[string]interface{}{}}
	for key, val := range p.Content {
		obj[key] = val
	}

	// Send the JSON object
	sendObj := []map[string]interface{}{obj}
	return c.ws.WriteJSON(sendObj)
}

func (c *Connection) WriteAck(p *Packet, ack interface{}) error {
	obj := map[string]interface{}{"id": p.Id, "channel": p.Channel,
		"ext": map[string]interface{}{"ack": ack}}
	for key, val := range p.Content {
		obj[key] = val
	}

	// Send the JSON object
	sendObj := []map[string]interface{}{obj}
	return c.ws.WriteJSON(sendObj)
}

func (c *Connection) readLoop() {
	for {
		if p, err := c.readRaw(); err != nil {
			// We have reached the end of the stream
			c.incoming <- incomingData{nil, err}
			close(c.incoming)
			return
		} else if p.Channel == "/meta/connect" {
			// Send a connect packet
			ext := p.Content["ext"].(map[string]interface{})
			content := map[string]interface{}{"clientId": c.clientId,
				"connectionType": "websocket"}
			p := c.Packet("/meta/connect", content)
			if err := c.WriteAck(p, ext["ack"]); err != nil {
				c.incoming <- incomingData{nil, err}
				close(c.incoming)
				return
			}
		} else {
			// Normal packet
			c.incoming <- incomingData{p, nil}
		}
	}
}

func (c *Connection) readRaw() (*Packet, error) {
	// Read the object
	var container []map[string]interface{}
	if err := c.ws.ReadJSON(&container); err != nil {
		return nil, err
	}
	if len(container) == 0 {
		return nil, errors.New("Got empty response")
	}
	object := container[0]

	// Decode the packet
	p := new(Packet)
	p.Content = map[string]interface{}{}

	// Decode channel and id keys
	if channel, ok := object["channel"]; !ok {
		return nil, errors.New("No 'channel' key")
	} else if p.Channel, ok = channel.(string); !ok {
		return nil, errors.New("Invalid type for 'channel' key")
	}
	if id, ok := object["id"]; ok {
		if p.Id, ok = id.(string); !ok {
			return nil, errors.New("Invalid type for 'id' key")
		}
	}

	for key, val := range object {
		if key == "id" || key == "channel" {
			continue
		}
		p.Content[key] = val
	}
	return p, nil
}
