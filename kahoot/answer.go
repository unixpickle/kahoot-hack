package kahoot

func (c *Connection) WaitQuestion() (*Packet, error) {
	for {
		p, err := c.ReadChannel("/service/player")
		if err != nil {
			return nil, err
		}
		if data, ok := p.Content["data"]; ok {
			m := data.(map[string]interface{})
			if m["type"].(string) == "message" {
				if _, ok := m["questionIndex"]; ok {
					return p, nil
				}
			}
		}
	}
}

func (c *Connection) SendAnswer(packet *Packet, choice int) {
	data := packet.Content["data"].(map[string]interface{})
	id := data["id"].(int)
	meta := map[string]interface{}{"lag": 13, "device": "hey/1.0"}
	content := map[string]interface{}{"choice": choice, "meta": meta}
	data = map[string]interface{}{"type": "message", "gameid": c.gameid,
		"host": "kahoot.it", "content": content, "id": id}
	c.WriteData("/service/controller", data)
}

