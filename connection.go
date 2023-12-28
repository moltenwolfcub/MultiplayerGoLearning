package main

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

func (c Connection) Send(packet Packet) {
	err := c.Encoder.Encode(&packet)
	if err != nil {
		log.Fatal("encoding error: ", err)
	}
}

func (c Connection) Recieve() Packet {
	var packet Packet

	err := c.Decoder.Decode(&packet)
	if err != nil {
		log.Fatal("decoding error: ", err)
	}
	return packet
}
