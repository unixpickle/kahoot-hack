package kahoot

func (c *Connection) ReadChannel(channel string) (*Packet, error) {
	return c.ReadFilter(func(p *Packet) bool {
		return p.Channel == channel
	})
}

func (c *Connection) ReadFilter(f PacketFilter) (*Packet, error) {
	for {
		p, err := c.Read()
		if err != nil {
			return nil, err
		}
		if f(p) {
			return p, nil
		}
	}
}

func (c *Connection) ReadId(id string) (*Packet, error) {
	return c.ReadFilter(func(p *Packet) bool {
		return p.Id == id
	})
}
