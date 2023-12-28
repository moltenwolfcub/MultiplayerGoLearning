package main

import (
	"fmt"
	"net"
)

type Client struct {
	listenAddr string
	connection Connection
}

func NewClient(listenAddr string) *Client {
	return &Client{
		listenAddr: listenAddr,
	}
}

func (c *Client) Start() error {
	conn, err := net.Dial("tcp", c.listenAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	c.connection = NewConnection(conn)

	for {
		rawPacket := c.connection.Recieve()

		switch packet := rawPacket.(type) {
		case ClientboundMessagePacket:
			fmt.Println(packet.Message)
		}
	}
}
