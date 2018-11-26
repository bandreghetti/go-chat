package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/bandreghetti/go-chat/msgs"
)

const (
	serverAddr     = "chat-server"
	serverPort     = "6174"
	loginMenu      = "Please choose an username:"
	welcomeMsg     = "Welcome, %s!\n"
	wrongCmdMsg    = "Invalid command. Enter '\\help' to list available commands"
	notInRoomMsg   = "You can't send messages if you're not in a room!\nEnter '\\help' to list available commands"
	joinInRoomMsg  = "You can't join a room without leaving the one you're in!"
	leaveNoRoomMsg = "You can't leave a room if you're not in one!"
)

var helpMsg = strings.Join([]string{
	"Available commands:",
	"\\list [room-name] - list available rooms or users in a room",
	"\\join <room-name> - join an existing room",
	"\\leave - leave the current room",
	"\\create <room-name> - create a new room",
	"\\delete <room-name> - delete an empty room",
	"\\logout - terminate the program",
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
		if response.Status == msgs.StatusInvalidUsername {
			fmt.Println("Invalid username!")
		} else if response.Status == msgs.StatusUsernameExists {
			fmt.Println("Username already exists! Please choose another one.")
		}
	}

	fmt.Printf(welcomeMsg, username)
	fmt.Println(helpMsg)

	inRoom := false
	loggedIn := true
	for loggedIn {
		scanner.Scan()
		command := scanner.Text()

		args := strings.Split(command, " ")
		if command[0] == '\\' {
			switch args[0] {
			case "\\list":
				if len(args) == 1 {
					request = msgs.ChatMsg{
						Command: msgs.CmdList,
					}
				} else if len(args) == 2 {
					request = msgs.ChatMsg{
						Command: msgs.CmdListUsers,
						Payload: []byte(args[1]),
					}
				}
				response = sendMsg(request)
				if response.Status != msgs.StatusOK {
					if response.Status == msgs.StatusInexistentRoom {
						fmt.Printf("Room named %s does not exist\n", args[1])
					}
					continue
				}
				fmt.Println(string(response.Payload))
			case "\\join":
				if len(args) != 2 {
					fmt.Println("Command \\join requires an argument")
					continue
				}
				if inRoom {
					fmt.Println(joinInRoomMsg)
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
				inRoom = true
				go updateMessages(&inRoom)
				fmt.Printf("Welcome to room %s!\n", args[1])
			case "\\leave":
				if !inRoom {
					fmt.Println(leaveNoRoomMsg)
					continue
				}
				inRoom = false
				request = msgs.ChatMsg{
					Command: msgs.CmdLeave,
				}
				sendMsg(request)
				fmt.Println("You left the room.")
			case "\\create":
				if inRoom {
					fmt.Println("Cannot manage rooms while inside one.")
					continue
				}
				if len(args) != 2 {
					fmt.Println("Command \\create requires an argument.")
					continue
				}
				request = msgs.ChatMsg{
					Command: msgs.CmdCreateRoom,
					Payload: []byte(args[1]),
				}
				response = sendMsg(request)
				if response.Status != msgs.StatusOK {
					if response.Status == msgs.StatusRoomAlreadyExists {
						fmt.Printf("Room named %s already exists.\n", args[1])
					}
					continue
				}
				fmt.Printf("Successfully created room named %s.\n", args[1])
			case "\\delete":
				if inRoom {
					fmt.Println("Cannot manage rooms while inside one.")
					continue
				}
				if len(args) != 2 {
					fmt.Println("Command \\delete requires an argument")
					continue
				}
				request = msgs.ChatMsg{
					Command: msgs.CmdDeleteRoom,
					Payload: []byte(args[1]),
				}
				response = sendMsg(request)
				if response.Status != msgs.StatusOK {
					if response.Status == msgs.StatusRoomNotEmpty {
						fmt.Println("Cannot delete a non-empty room!")
					} else if response.Status == msgs.StatusInexistentRoom {
						fmt.Printf("Room named %s does not exist.\n", args[1])
					}

					continue
				}
				fmt.Printf("Successfully deleted room named %s.\n", args[1])
			case "\\logout":
				request = msgs.ChatMsg{
					Command: msgs.CmdLogout,
				}
				inRoom = false
				sendMsg(request)
				loggedIn = false
			case "\\help":
				fmt.Println(helpMsg)
			default:
				fmt.Println(wrongCmdMsg)
			}
		} else {
			if inRoom {
				request = msgs.ChatMsg{
					Command: msgs.CmdMsg,
					Payload: []byte(command),
				}
				response = sendMsg(request)
				if response.Status != msgs.StatusOK {
					fmt.Printf("Couldn't post message")
				}
			} else {
				fmt.Println(notInRoomMsg)
			}
		}
	}
}

func updateMessages(inRoom *bool) {
	// Get index from which client should request messages
	request := msgs.ChatMsg{
		Command: msgs.CmdGetMsgIdx,
	}
	response := sendMsg(request)
	lastIdx := response.Payload

	// Fetch new messages every second
	tick := time.Tick(500 * time.Millisecond)
	for *inRoom {
		_ = <-tick
		request := msgs.ChatMsg{
			Command: msgs.CmdFetch,
			Payload: lastIdx,
		}
		response = sendMsg(request)
		if response.Status != msgs.StatusOK {
			// TODO: handle errors
		}
		messages := string(response.Payload[:len(response.Payload)-8])
		lastIdx = response.Payload[len(response.Payload)-8:]
		if len(messages) > 0 {
			fmt.Println(messages)
		}
	}
}
