package server

import (
	"bufio"
	"fmt"
	"net"
)

type Client chan<- string

func Broadcaster(entering, leaving chan Client, messages chan string) {
	clients := make(map[Client]bool)
	for {
		select {
		case msg := <-messages:
			for c := range clients {
				c <- msg
			}
		case c := <-entering:
			clients[c] = true
		case c := <-leaving:
			delete(clients, c)
			close(c)
		}
	}
}

func HandleConn(conn net.Conn, entering, leaving chan Client, messages chan string) {
	ch := make(chan string)
	go ClientWriter(conn, ch)
	who := conn.RemoteAddr().String()
	ch <- "You are " + who
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

func ClientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
