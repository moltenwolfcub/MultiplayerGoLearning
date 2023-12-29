package common

import (
	"encoding/gob"
	"net"
)

type Packet interface {
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
	Message string
}

type ServerboundAnnouncePacket struct {
	Announcement string
}
