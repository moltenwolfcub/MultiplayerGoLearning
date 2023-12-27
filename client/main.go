package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":2525")
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	defer conn.Close()

	conn.Write([]byte("test"))
}
