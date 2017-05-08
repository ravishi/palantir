package bennu

import "github.com/gorilla/websocket"

type connection struct {
	b *broadcaster
	c *websocket.Conn
}


func createConnection(conn *websocket.Conn) *connection {
	return &connection{
		b: createBroadcaster(),
		c: conn,
	}
}

func (c *connection) Handle() error {
	return nil
}

func (c *connection) Close(upErr error) error {
	c.b.Close()

	return upErr
}
