package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitCh     chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitCh:     make(chan struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitCh

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		fmt.Println("New connection to the server: ", conn.RemoteAddr())

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
			return
		}

		if err != nil {
			fmt.Println("read error: ", err)
			continue
		}

		msg := buf[:n]
		fmt.Println(string(msg))
		conn.Write(msg)
	}
}

func main() {
	server := NewServer(":2525")
	log.Fatal(server.Start())
}
