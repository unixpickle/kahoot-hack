package kahoot

import "errors"
import "fmt"

func (c *Connection) Subscribe(channels []string) error {
	for _, ch := range channels {
		content := map[string]interface{}{"subscription": ch,
			"clientId": c.clientId}
		pack := c.NewPacket("/meta/subscribe", content)
		if err := c.Write(pack, false); err != nil {
			return err
		}
		res, err := c.Read()
		if err != nil {
			return err
		}
		// Check the 'success' field
		fmt.Println("channel:", res.Channel, "content:", res.Content)
		success, ok := res.Content["successful"].(bool)
		if !ok || !success {
			return errors.New("Negative 'successful' field for channel: " + ch)
		}
	}
	return nil
}
