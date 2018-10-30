package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

const (
	serverPort = ":6174"
	connType   = "tcp"
)

func main() {
	// Listen on TCP port serverPort.
	l, err := net.Listen(connType, serverPort)
	if err != nil {
		log.Println("error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener object after returning.
	defer l.Close()
	log.Println("Listening on " + l.Addr().String())

	// Loop over TCP requests
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Close the connection when you're done with it.
	defer conn.Close()

	// Read message received
	buf := bufio.NewReader(conn)
	msg, err := buf.ReadString('\n')
	if err != nil {
		log.Println("Error reading:", err.Error())
		return
	}

	log.Printf("Received %d bytes from %s", len(msg), conn.RemoteAddr().String())
	log.Printf("Message received: %s\n", msg)

	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
}
