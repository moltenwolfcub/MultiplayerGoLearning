package client

import (
	"fmt"
	"log"
	"net"

	"github.com/moltenwolfcub/MultiplayerGoLearning/common"
)

type Client struct {
	listenAddr string
	connection common.Connection
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

	c.connection = common.NewConnection(conn)

	go c.readLoop()
	return c.mainLoop()
}

func (c *Client) readLoop() error {
	for {
		rawPacket := c.connection.MustRecieve()
		err := c.handlePacket(rawPacket)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func (c *Client) mainLoop() error {
	for {
		var message string
		fmt.Print(">>> ")
		fmt.Scanln(&message)
		c.connection.MustSend(common.ServerboundAnnouncePacket{Announcement: message})
	}
}

func (c *Client) handlePacket(rawPacket common.Packet) error {
	switch packet := rawPacket.(type) {
	case common.ClientboundMessagePacket:
		fmt.Print("\033[2K\r" + packet.Message + "\n>>> ")
	default:
		return fmt.Errorf("unkown packet: %s", packet)
	}
	return nil
}
