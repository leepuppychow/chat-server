/* EXERCISES:

8.13:
	Make the server disconnect idle clients. Hint: calling conn.Close() in another goroutine
	unblocks active Read calls such as the one done by input.Scan()

*/

package main

import (
	"fmt"
	"log"
	"net"

	s "github.com/leepuppychow/chat-server/server"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	fmt.Println("Chat server running on port 8080")
	if err != nil {
		log.Fatal(err)
	}
	go s.Broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.HandleConn(conn)
	}
}
