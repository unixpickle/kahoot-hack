package kahoot

import (
	"encoding/json"
)

func (c *Connection) WaitQuestion() (*Packet, error) {
	for {
		p, err := c.ReadChannel("/service/player")
		if err != nil {
			return nil, err
		}
		if data, ok := p.Content["data"]; ok {
			m := data.(map[string]interface{})
			if content, ok := m["content"]; ok {
				m := map[string]interface{}{}
				err := json.Unmarshal([]byte(content.(string)), &m)
				if err == nil {
					if _, ok := m["questionNumber"]; ok {
						return p, nil
					}
				}
			}
		}
	}
}

func (c *Connection) SendAnswer(packet *Packet, choice int) error {
	//data := packet.Content["data"].(map[string]interface{})
	//id := data["id"]
	screen := map[string]interface{}{"width": 1920, "height": 1080}
	device := map[string]interface{}{"userAgent": "hey", "screen": screen}
	meta := map[string]interface{}{"lag": 100000000, "device": device}
	content := map[string]interface{}{"choice": choice, "meta": meta}
	if enc, err := json.Marshal(content); err != nil {
		return err
	} else {
		data := map[string]interface{}{"type": "message", "gameid": c.gameid,
			"host": "kahoot.it", "content": string(enc), "id": 6}
		_, err = c.WriteData("/service/controller", data)
		return err
	}
}

