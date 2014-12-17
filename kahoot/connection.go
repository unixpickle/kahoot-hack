package kahoot

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
)

type Connection struct {
	ws *websocket.Conn
}

type ClientPacket struct {
	
}

func NewConnection(gamePin string) (*Connection, error) {
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
	return &Connection{ws}, nil
}
