package Client

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	Conn *websocket.Conn
	roomCommandChannel chan string
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{Conn:conn}
}

func (c *Client) SetRoomCommandChannel(commandChannel chan string) {
	c.roomCommandChannel = commandChannel
}

func (c *Client) Read() {
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Client error: %v", err)
			break
		}

		messageAsString := string(message)
		c.roomCommandChannel <- messageAsString
	}
}