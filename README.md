## Learning repo: Chat Server in Go

* Resources used: *The Go Programming Language* (Donovan and Kernighan, 2016)

## To run chat server:

1. From root directory: `go run main.go` 
2. Open as many terminal shells as you would like.
3. Connect to chat server with: `telnet 127.0.0.1 8080`

## NOTES:

When client connects, HandleConn and ClientWriter goroutines are started. Also a string channel (line 12 in connection) is created and sent to ClientWriter. This same channel is sent to the global entering channel. 

The Broadcaster goroutine, which has already been running, will receive from the entering channel. It will add that same channel (which is the Client) to the clients map. 

The ClientWriter goroutine is running in an endless loop because the for range loop for a channel will keep listening unless you close that channel.

Now, whenever a string is sent via the global messages channel, the Broadcaster goroutine takes that string and sends it to each client (channel) within the clients map. 

```
  case msg := <-messages:
    for c := range clients {
      c <- msg
    }
```

This client (channel) is the same one that the ClientWriter goroutine is listening to (for each respective client). This is how the message gets printed to each client's stdout. 

