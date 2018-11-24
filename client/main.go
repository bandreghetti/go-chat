package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/bandreghetti/go-chat/msgs"
)

const (
	serverAddr = "chat-server"
	serverPort = "6174"
	loginMenu  = "Please choose an username:"
	welcomeMsg = "Welcome, %s!\n"
)

func sendMsg(message msgs.ChatMsg) msgs.ChatMsg {
	msgJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error generating msgJSON: %s", err.Error())
		os.Exit(1)
	}
	// Connect to server
	conn, err := net.Dial("tcp", serverAddr+":"+serverPort)
	if err != nil {
		log.Printf("Error dialing server: %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// Write to socket
	conn.Write(msgJSON)

	// Listen for reply
	var recvMsg msgs.ChatMsg
	json.NewDecoder(conn).Decode(&recvMsg)
	return recvMsg
}

func main() {
	var username string
	scanner := bufio.NewScanner(os.Stdin)
	var response msgs.ChatMsg
	for response.Status != msgs.StatusOK {
		fmt.Println(loginMenu)
		scanner.Scan()
		username = scanner.Text()

		login := msgs.ChatMsg{
			Command: msgs.CmdLogin,
			Payload: []byte(username),
		}
		response = sendMsg(login)
	}

	fmt.Printf(welcomeMsg, username)
}
