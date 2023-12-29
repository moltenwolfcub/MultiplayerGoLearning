package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/moltenwolfcub/MultiplayerGoLearning/client"
	"github.com/moltenwolfcub/MultiplayerGoLearning/common"
	"github.com/moltenwolfcub/MultiplayerGoLearning/server"
)

var sideFlag string

func main() {
	common.RegisterPackets()

	flag.StringVar(&sideFlag, "side", "", "'server' or 'client'")
	flag.Parse()

	if sideFlag == "server" {
		fmt.Println("server")

		server := server.NewServer(":2525")

		log.Fatal(server.Start())

	} else if sideFlag == "client" {
		fmt.Println("client")
		client := client.NewClient(":2525")

		log.Fatal(client.Start())
	} else {
		log.Fatal("Unknown side to launch")
	}
}
