package main

import (
	"fmt"
	"log"
	"net"

	"github.com/leepuppychow/chat-server/clients"
	"github.com/leepuppychow/chat-server/server"
)

var (
	entering = make(chan clients.Client)
	leaving  = make(chan clients.Client)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	fmt.Println("Chat server running on port 8080")
	if err != nil {
		log.Fatal(err)
	}
	go server.Broadcaster(entering, leaving, messages)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go clients.HandleConn(conn, entering, leaving, messages)
	}
}
