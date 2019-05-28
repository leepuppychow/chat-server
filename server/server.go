package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	name string
	out  chan<- string
}

func Broadcaster(entering, leaving chan Client, messages chan string) {
	clients := make(map[Client]bool)
	for {
		select {
		case msg := <-messages:
			for c := range clients {
				c.out <- msg
			}
		case c := <-entering:
			clients[c] = true
			displayAll(clients)
		case c := <-leaving:
			delete(clients, c)
			close(c.out)
			displayAll(clients)
		}
	}
}

func displayAll(clients map[Client]bool) {
	all := "All current clients:\n\n"
	for c := range clients {
		all += "\t" + c.name + "\n"
	}
	for c := range clients {
		c.out <- all
	}
}

func HandleConn(conn net.Conn, entering, leaving chan Client, messages chan string) {
	ch := make(chan string)
	go ClientWriter(conn, ch)

	reader := bufio.NewReader(conn)
	ch <- "Enter Username: "
	name, _ := reader.ReadString('\n')

	newClient := Client{
		out:  ch,
		name: strings.TrimSpace(name),
	}
	entering <- newClient
	messages <- newClient.name + " has arrived"

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- newClient.name + ": " + input.Text()
	}

	leaving <- newClient
	messages <- newClient.name + " has left."
	conn.Close()
}

func ClientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
