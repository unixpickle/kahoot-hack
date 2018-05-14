package kahoot

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestGameSessionToken(t *testing.T) {
	resp, err := http.Post("https://play.kahoot.it/reserve/session/", "text/plain",
		bytes.NewReader(nil))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		t.Fatal("reserve session:", err)
	}
	pin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("read session pin:", err)
	}

	token, err := gameSessionToken(strings.TrimSpace(string(pin)))
	if err != nil {
		t.Fatal("get token:", err)
	}

	// The server doesn't let us establish a WebSocket unless
	// we present a valid session token.
	conn, err := net.Dial("tcp", "kahoot.it:443")
	if err != nil {
		t.Fatal("website out of reach:", err)
	}
	defer conn.Close()

	url, err := url.Parse("wss://kahoot.it/cometd/" + string(pin) + "/" + token)
	if err != nil {
		t.Fatal(err)
	}
	reqHeader := http.Header{}
	reqHeader.Set("Origin", "https://kahoot.it")
	reqHeader.Set("Cookie", "no.mobitroll.session="+string(pin))
	ws, _, err := websocket.NewClient(conn, url, reqHeader, 100, 100)
	defer ws.Close()
	if err != nil {
		t.Fatal("establish WebSocket:", err)
	}
}
