package kahoot

func (c *Connection) WriteData(channel string,
	data interface{}) (string, error) {
	content := map[string]interface{}{"data": data, "clientId": c.clientId}
	pack := c.NewPacket(channel, content)
	if err := c.Write(pack, false); err != nil {
		return "", err
	}
	return pack.Id, nil
}
