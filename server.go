package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type Server struct {
	listenAddr string
	listener   net.Listener
	quitCh     chan struct{}
	inMsgCh    chan Packet
	peers      map[net.Addr]Connection
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitCh:     make(chan struct{}),
		inMsgCh:    make(chan Packet, 10),
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
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)

		if errors.Is(err, io.EOF) {
			fmt.Println("Lost connection to peer: ", conn.RemoteAddr())
			delete(s.peers, conn.RemoteAddr())
			return
		}

		if err != nil {
			fmt.Println("read error: ", err)
			continue
		}

		fmt.Println("Recieved from connection:", buf[:n])
	}
}

func (s *Server) mainLoop() {
	for {
		time.Sleep(time.Second * 3)

		for _, conn := range s.peers {
			conn.Send(ClientboundMessagePacket{Message: "hello"})
		}
	}
}
