package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var (
	entering = make(chan Client)
	leaving  = make(chan Client)
	messages = make(chan string)
)

type Client struct {
	name string
	out  chan<- string
}

func Broadcaster() {
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

func HandleConn(conn net.Conn) {
	ch := make(chan string)
	go ClientWriter(conn, ch)
	cli := NewClient(conn, ch)

	cli.SendMessage(conn)
	cli.Leaving(conn)
}

func ClientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func NewClient(conn net.Conn, ch chan string) Client {
	reader := bufio.NewReader(conn)
	ch <- "Enter Username: "
	name, _ := reader.ReadString('\n') // Ignoring errors here
	newClient := Client{
		out:  ch,
		name: strings.TrimSpace(name),
	}
	entering <- newClient
	messages <- newClient.name + " has arrived"
	return newClient
}

func (c *Client) SendMessage(conn net.Conn) {
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- c.name + ": " + input.Text()
	}
}

func (c *Client) Leaving(conn net.Conn) {
	leaving <- *c
	messages <- c.name + " has left."
	conn.Close()
}
