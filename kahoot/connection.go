package kahoot

import (
	"crypto/tls"
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"strconv"
)

type PacketFilter func(*Packet) bool

type Connection struct {
	ws       *websocket.Conn
	id       int
	gameid   int
	clientId string
}

type Packet struct {
	Channel string
	Id      string
	Content map[string]interface{}
}

func NewConnection(gamePin string) (*Connection, error) {
	gameid, err := strconv.Atoi(gamePin)
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
	reqHeader.Set("Cookie", "no.mobitroll.session=" + gamePin)
	ws, _, err := websocket.NewClient(conn, url, reqHeader, 100, 100)
	if err != nil {
		return nil, err
	}
	return &Connection{ws, 0, gameid, ""}, nil
}

func (c *Connection) NewPacket(channel string,
	content map[string]interface{}) *Packet {
	c.id++
	return &Packet{channel, strconv.Itoa(c.id), content}
}

func (c *Connection) Write(p *Packet, ack interface{}) error {
	obj := map[string]interface{}{"id": p.Id, "channel": p.Channel}
	for key, val := range p.Content {
		obj[key] = val
	}
	ext := map[string]interface{}{}
	if b, ok := ack.(bool); !ok || b {
		ext["ack"] = ack
	}
	obj["ext"] = ext
	
	// Send the JSON object
	sendObj := []map[string]interface{}{obj}
	return c.ws.WriteJSON(sendObj)
}

func (c *Connection) Read() (*Packet, error) {
	for {
		p, err := c.readRaw()
		if err != nil {
			return nil, err
		}
		if p.Channel == "/meta/connect" {
			// Reply to it
			ext := p.Content["ext"].(map[string]interface{})
			c.SendConnect(ext["ack"])
		} else {
			return p, err
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

func (c *Connection) ReadId(id string) (*Packet, error) {
	return c.readFilter(func(p *Packet) bool {
		return p.Id == id
	})
}

func (c *Connection) ReadChannel(channel string) (*Packet, error) {
	return c.readFilter(func(p *Packet) bool {
		return p.Channel == channel
	})
}

func (c *Connection) readFilter(f PacketFilter) (*Packet, error) {
	for {
		p, err := c.Read()
		if err != nil {
			return nil, err
		}
		if f(p) {
			return p, nil
		}
		// Discarding packet
	}
}
