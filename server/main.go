package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/bandreghetti/go-chat/msgs"
)

const (
	serverPort = "chat-server:6174"
	connType   = "tcp"
)

var ip2user = map[string]string{
	"127.0.0.1": "server",
}

var user2ip = map[string]string{
	"server": "127.0.0.1",
}

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
			log.Println("Error accepting connection: ", err.Error())
		} else {
			go handleRequest(conn)
		}
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	time.Sleep(1000 * time.Millisecond)

	// Close the connection when you're done with it.
	defer conn.Close()

	// Read message received
	var recvMsg msgs.ChatMsg
	err := json.NewDecoder(conn).Decode(&recvMsg)
	if err != nil {
		log.Println("Invalid JSON received: " + err.Error())
		return
	}

	switch recvMsg.Command {
	case msgs.CmdLogin:
		login(conn, recvMsg)
	}
}

func login(conn net.Conn, recvMsg msgs.ChatMsg) {
	// Get requesting IP
	requestAddr := strings.Split(conn.RemoteAddr().String(), ":")
	requestIP := requestAddr[0]

	username := string(recvMsg.Payload)

	if !validUsername(username) {
		// respond with error
	}

	var response msgs.ChatMsg
	ip2user[requestIP] = username
	user2ip[username] = requestIP
	response = msgs.ChatMsg{
		Status: msgs.StatusOK,
	}
	respond(conn, response)
	log.Printf("%s requested login as %s", requestIP, username)
}

func respond(conn net.Conn, response msgs.ChatMsg) {
	respJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("error marshaling response: %s", err.Error())
	}
	conn.Write(respJSON)
}
