package main

import (
	"errors"
	"fmt"
	"io"
	"net"
)

type Server struct {
	listenAddr string
	listener   net.Listener
	quitCh     chan struct{}
	inMsgCh    chan RecievedPacket
	peers      map[net.Addr]Connection
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitCh:     make(chan struct{}),
		inMsgCh:    make(chan RecievedPacket, 10),
		peers:      make(map[net.Addr]Connection),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	s.listener = listener

	go s.mainLoop()
	go s.acceptLoop()

	<-s.quitCh
	close(s.inMsgCh)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		fmt.Println("New connection to the server: ", conn.RemoteAddr())
		s.peers[conn.RemoteAddr()] = NewConnection(conn)

		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	for {
		rawPacket, err := s.peers[conn.RemoteAddr()].Recieve()

		if errors.Is(err, io.EOF) {
			fmt.Println("Lost connection to peer: ", conn.RemoteAddr())
			delete(s.peers, conn.RemoteAddr())
			return
		}

		if err != nil {
			fmt.Println("read error: ", err.Error())
			continue
		}

		s.inMsgCh <- RecievedPacket{
			Packet: rawPacket,
			Sender: conn.RemoteAddr(),
		}
	}
}

func (s *Server) mainLoop() {
	for {
		for rawPacket := range s.inMsgCh {
			s.handlePacket(rawPacket)
		}
	}
}

func (s *Server) handlePacket(recieved RecievedPacket) error {
	switch packet := recieved.Packet.(type) {
	case ServerboundAnnouncePacket:
		s.announce(packet.Announcement, recieved.Sender)
	default:
		return fmt.Errorf("unkown packet: %s", packet)
	}
	return nil
}

func (s *Server) announce(announcement string, sender net.Addr) {
	fmt.Printf("Connection %v sent announcement with message: %s\n", sender, announcement)
	for addr, conn := range s.peers {
		if addr == sender {
			conn.MustSend(ClientboundMessagePacket{Message: "Your announcment of: '" + announcement + "' has been sent to everyone"})
			continue
		}

		conn.MustSend(ClientboundMessagePacket{Message: announcement})
	}
}
