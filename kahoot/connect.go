package kahoot

func (c *Connection) SendConnect(ackNum interface{}) error {
	// Send connect packet
	content := map[string]interface{}{"clientId": c.clientId,
		"connectionType": "websocket"}
	p := c.NewPacket("/meta/connect", content)
	if err := c.Write(p, ackNum); err != nil {
		return err
	}
	return nil
}
