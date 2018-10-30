package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	serverAddr = "chat-server"
	serverPort = "6174"
)

func main() {
	for i := 0; ; i++ {
		time.Sleep(time.Second)
		// Connect to server
		conn, err := net.Dial("tcp", serverAddr+":"+serverPort)
		if err != nil {
			log.Printf("error dialing server: %s", err.Error())
			continue
		}
		// send to socket
		fmt.Fprintf(conn, "test-message %d\n", i)
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
