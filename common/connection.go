package common

import (
	"encoding/gob"
	"log"
	"net"
)

type Connection struct {
	Connection net.Conn
	ClientAddr net.Addr
	Encoder    *gob.Encoder
	Decoder    *gob.Decoder
}

func NewConnection(conn net.Conn) Connection {
	return Connection{
		Connection: conn,
		ClientAddr: conn.RemoteAddr(),
		Encoder:    gob.NewEncoder(conn),
		Decoder:    gob.NewDecoder(conn),
	}
}

/*
Encodes a packet with the gob.Encoder and sends it
along the Connection
*/
func (c Connection) Send(packet Packet) error {
	err := c.Encoder.Encode(&packet)
	if err != nil {
		return err
	}
	return nil
}

/*
Waits for incoming data on the Connection before it
decodes the data into a packet with the gob.Decoder
*/
func (c Connection) Recieve() (Packet, error) {
	var packet Packet

	err := c.Decoder.Decode(&packet)
	if err != nil {
		return nil, err
	}
	return packet, nil
}

/*
Same as Send() but panics if it encounters an error
*/
func (c Connection) MustSend(packet Packet) {
	err := c.Encoder.Encode(&packet)
	if err != nil {
		log.Fatal("encoding error: ", err)
	}
}

/*
Same as Recieve() but panics if it encounters an error
*/
func (c Connection) MustRecieve() Packet {
	var packet Packet

	err := c.Decoder.Decode(&packet)
	if err != nil {
		log.Fatal("decoding error: ", err)
	}
	return packet
}
