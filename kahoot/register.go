package kahoot

func (c *Connection) Register(nick string) error {
	// Handshake, subscribe, connect, register
	if err := c.Handshake(); err != nil {
		return err
	}
	if err := c.SubscribeAll(); err != nil {
		return err
	}
	if err := c.SendConnect(-1); err != nil {
		return err
	}

	if err := c.UnsubscribeAll(); err != nil {
		return err
	}
	if err := c.SubscribeAll(); err != nil {
		return err
	}
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

func (c *Connection) SubscribeAll() error {
	return c.Subscribe("/service/player", "/service/controller",
		"/service/status")
}

func (c *Connection) UnsubscribeAll() error {
	return c.Unsubscribe("/service/player", "/service/controller",
		"/service/status")
}
