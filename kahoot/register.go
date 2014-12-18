package kahoot

func (c *Connection) Register(nick string) error {
	// Handshake, subscribe, connect, register
	if err := c.Handshake(); err != nil {
		return err
	}
	c.Subscribe("/service/player", "/service/controller", "/service/status")
	loginData := map[string]interface{}{"type": "login", "gameid": c.gameid,
		"host": "kahoot.it", "name": nick}
	lastId, err := c.WriteData("/service/controller", loginData)
	if err != nil {
		return err
	}
	if _, err := c.ReadId(lastId); err != nil {
		return err
	}
	return nil
}