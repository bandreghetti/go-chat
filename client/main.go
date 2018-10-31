package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"./validation"
)

const (
	serverAddr = "chat-server"
	serverPort = "6174"
	loginMenu  = "Please choose an username:"
	welcomeMsg = "Welcome, %s!\n"
)

func main() {
	var username string
	scanner := bufio.NewScanner(os.Stdin)
	for !validation.ValidUsername(username) {
		fmt.Println(loginMenu)
		scanner.Scan()
		username = scanner.Text()
	}

	fmt.Printf(welcomeMsg, username)

	for {
		var msg string
		fmt.Scanln(&msg)
		// Connect to server
		conn, err := net.Dial("tcp", serverAddr+":"+serverPort)
		if err != nil {
			log.Printf("error dialing server: %s", err.Error())
			continue
		}
		// send to socket
		fmt.Fprintf(conn, "%s\n", msg)
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
