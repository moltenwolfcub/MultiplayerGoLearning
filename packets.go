package main

import (
	"encoding/gob"
	"net"
)

type Packet interface {
	markPacket()
}

type RecievedPacket struct {
	Packet Packet
	Sender net.Addr
}

func RegisterPackets() {
	gob.Register(ClientboundMessagePacket{})

	gob.Register(ServerboundAnnouncePacket{})
}

type ClientboundMessagePacket struct {
	Packet
	Message string
}

type ServerboundAnnouncePacket struct {
	Packet
	Announcement string
}
