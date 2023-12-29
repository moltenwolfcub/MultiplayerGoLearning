package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/moltenwolfcub/MultiplayerGoLearning/common"
)

type Server struct {
	listenAddr string
	listener   net.Listener
	quitCh     chan struct{}
	inMsgCh    chan common.RecievedPacket
	peers      map[net.Addr]common.Connection
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitCh:     make(chan struct{}),
		inMsgCh:    make(chan common.RecievedPacket, 10),
		peers:      make(map[net.Addr]common.Connection),
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
	go s.packetLoop()
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
		s.peers[conn.RemoteAddr()] = common.NewConnection(conn)

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

		s.inMsgCh <- common.RecievedPacket{
			Packet: rawPacket,
			Sender: conn.RemoteAddr(),
		}
	}
}

func (s *Server) packetLoop() {
	for rawPacket := range s.inMsgCh {
		s.handlePacket(rawPacket)
	}
}

func (s *Server) mainLoop() {
	for {
		time.Sleep(time.Second)

		for _, conn := range s.peers {
			conn.MustSend(common.ClientboundMessagePacket{Message: "tick"})
		}
		time.Sleep(time.Second)

		for _, conn := range s.peers {
			conn.MustSend(common.ClientboundMessagePacket{Message: "tock"})
		}
	}
}

func (s *Server) handlePacket(recieved common.RecievedPacket) error {
	switch packet := recieved.Packet.(type) {
	case common.ServerboundAnnouncePacket:
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
			conn.MustSend(common.ClientboundMessagePacket{Message: "Your announcment of: '" + announcement + "' has been sent to everyone"})
			continue
		}

		conn.MustSend(common.ClientboundMessagePacket{Message: announcement})
	}
}
