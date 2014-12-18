package kahoot

func (c *Connection) WriteData(ch string, data interface{}) (string, error) {
	content := map[string]interface{}{"data": data, "clientId": c.clientId}
	pack := c.Packet(ch, content)
	if err := c.Write(pack); err != nil {
		return "", err
	}
	return pack.Id, nil
}
