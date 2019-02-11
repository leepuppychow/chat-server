package clients

import (
	"bufio"
	"net"
)

type Client chan<- string

func HandleConn(conn net.Conn, entering, leaving chan Client, messages chan string) {
	ch := make(chan string)
	go ClientWriter(conn, ch)
	who := conn.RemoteAddr().String()
	ch <- "You are" + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}

	leaving <- ch
	messages <- who + " has left."
	conn.Close()
}
