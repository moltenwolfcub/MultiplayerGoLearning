package main

import "encoding/gob"

type Packet interface {
	markPacket()
}

func RegisterPackets() {
	gob.Register(ClientboundMessagePacket{})
}

type ClientboundMessagePacket struct {
	Packet
	Message string
}
