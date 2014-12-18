package kahoot

import "errors"

func (c *Connection) Handshake() error {
	content := map[string]interface{}{"version": "1.0", "minimumVersion": "1.0",
		"supportedConnectionTypes": []string{"websocket"},
		"advice": map[string]int{"timeout": 60000,"interval": 0}}
	pack := c.NewPacket("/meta/handshake", content)
	if err := c.Write(pack, true); err != nil {
		return err
	}
	res, err := c.ReadId(pack.Id)
	if err != nil {
		return err
	}
	
	// Check the 'success' field
	success, ok := res.Content["successful"].(bool)
	if !ok || !success {
		return errors.New("Negative 'successful' field.")
	}
	
	// Read the 'clientId' field
	c.clientId, ok = res.Content["clientId"].(string)
	if !ok {
		return errors.New("No 'clientId' response field")
	}
	
	return nil
}
