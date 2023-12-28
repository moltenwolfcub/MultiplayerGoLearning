package main

import (
	"flag"
	"fmt"
	"log"
)

var sideFlag string

func main() {
	RegisterPackets()

	flag.StringVar(&sideFlag, "side", "", "'server' or 'client'")
	flag.Parse()

	if sideFlag == "server" {
		fmt.Println("server")

		server := NewServer(":2525")

		log.Fatal(server.Start())

	} else if sideFlag == "client" {
		fmt.Println("client")
		client := NewClient(":2525")

		log.Fatal(client.Start())
	} else {
		log.Fatal("Unknown side to launch")
	}
}
