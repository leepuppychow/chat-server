package server

import "github.com/leepuppychow/chat-server/clients"

func Broadcaster(entering, leaving chan clients.Client, messages chan string) {
	clients := make(map[clients.Client]bool)
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
