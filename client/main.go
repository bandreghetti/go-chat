package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/bandreghetti/go-chat/msgs"
)

const (
	serverAddr  = "chat-server"
	serverPort  = "6174"
	loginMenu   = "Please choose an username:"
	welcomeMsg  = "Welcome, %s!\n"
	wrongCmdMsg = "Invalid command. Enter '\\help' to list available commands"
)

var helpMsg = strings.Join([]string{
	"Available commands:",
	"\\list - list available rooms",
	"\\join - join an existing room",
	"\\help - list available commands",
}, "\n")

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
	var request msgs.ChatMsg
	var response msgs.ChatMsg

	var username string
	scanner := bufio.NewScanner(os.Stdin)
	for response.Status != msgs.StatusOK {
		fmt.Println(loginMenu)
		scanner.Scan()
		username = scanner.Text()

		request = msgs.ChatMsg{
			Command: msgs.CmdLogin,
			Payload: []byte(username),
		}
		response = sendMsg(request)
	}

	fmt.Printf(welcomeMsg, username)
	fmt.Println(helpMsg)

	// inRoom := false
	loggedIn := true
	for loggedIn {
		scanner.Scan()
		command := scanner.Text()

		args := strings.Split(command, " ")
		switch args[0] {
		case "\\list":
			request = msgs.ChatMsg{
				Command: msgs.CmdList,
			}
			response = sendMsg(request)
			fmt.Println(string(response.Payload))
		case "\\join":
			if len(args) != 2 {
				fmt.Println("Command \\join requires an argument")
				continue
			}
			request = msgs.ChatMsg{
				Command: msgs.CmdJoin,
				Payload: []byte(args[1]),
			}
			response = sendMsg(request)
			if response.Status != msgs.StatusOK {
				// TODO: handle errors
				continue
			}
			fmt.Printf("Welcome to room %s!\n", args[1])
		case "\\logout":
			// TODO: send logout request
			loggedIn = false
		case "\\help":
			fmt.Println(helpMsg)
		default:
			fmt.Println(wrongCmdMsg)
		}
	}
}
