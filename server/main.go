package main

import (
	"encoding/binary"
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

var (
	ip2user = map[string]string{
		"127.0.0.1": "server",
	}
	user2ip = map[string]string{
		"server": "127.0.0.1",
	}
	ip2roomName = map[string]string{}
	rooms       = map[string]*room{}
)

func main() {
	// Listen on TCP port serverPort.
	l, err := net.Listen(connType, serverPort)
	if err != nil {
		log.Println("error listening:", err.Error())
		os.Exit(1)
	}

	rooms["general"] = &room{
		roomName:  "general",
		userCount: 0,
	}
	rooms["general"].users = make(map[string]struct{})

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
	case msgs.CmdList:
		list(conn)
	case msgs.CmdJoin:
		join(conn, recvMsg)
	case msgs.CmdMsg:
		postMsg(conn, recvMsg)
	case msgs.CmdFetch:
		fetch(conn, recvMsg)
	case msgs.CmdGetMsgIdx:
		getMsgIdx(conn, recvMsg)
	case msgs.CmdLeave:
		leave(conn)
	case msgs.CmdCreateRoom:
		createRoom(conn, recvMsg)
	case msgs.CmdLogout:
		logout(conn)
	}
}

func login(conn net.Conn, recvMsg msgs.ChatMsg) {
	// Get requesting IP
	requestAddr := strings.Split(conn.RemoteAddr().String(), ":")
	requestIP := requestAddr[0]

	username := string(recvMsg.Payload)

	if !validUsername(username) {
		// TODO: respond with error
	}

	var response msgs.ChatMsg
	ip2user[requestIP] = username
	user2ip[username] = requestIP
	response = msgs.ChatMsg{
		Status: msgs.StatusOK,
	}
	respond(conn, response)
	log.Printf("%s logged in as %s", requestIP, username)
}

func list(conn net.Conn) {
	var roomSlice []string
	for _, room := range rooms {
		roomSlice = append(roomSlice, room.String())
	}
	payload := strings.Join(roomSlice, "\n")
	response := msgs.ChatMsg{
		Status:  msgs.StatusOK,
		Payload: []byte(payload),
	}
	respond(conn, response)
}

func join(conn net.Conn, recvMsg msgs.ChatMsg) {
	// Get requesting IP
	requestAddr := strings.Split(conn.RemoteAddr().String(), ":")
	requestIP := requestAddr[0]

	var response msgs.ChatMsg
	roomName := string(recvMsg.Payload)
	_, roomExists := rooms[roomName]
	if !roomExists {
		response.Status = msgs.StatusInexistentRoom
		respond(conn, response)
		return
	}
	rooms[roomName].JoinUser(requestIP)
	ip2roomName[requestIP] = roomName
	response.Status = msgs.StatusOK
	respond(conn, response)
	username := ip2user[requestIP]
	log.Printf("%s (%s) joined %s\n", username, requestIP, roomName)
}

func postMsg(conn net.Conn, recvMsg msgs.ChatMsg) {
	// Get requesting IP
	requestAddr := strings.Split(conn.RemoteAddr().String(), ":")
	requestIP := requestAddr[0]

	roomName, ok := ip2roomName[requestIP]
	if !ok {
		// TODO: handle error when ip is not in a room
	}

	r := rooms[roomName]

	status := r.PostMessage(requestIP, string(recvMsg.Payload))
	response := msgs.ChatMsg{
		Status: status,
	}
	respond(conn, response)
	username := ip2user[requestIP]
	log.Printf("%s (%s) posted %s\n", username, requestIP, string(recvMsg.Payload))
}

func fetch(conn net.Conn, recvMsg msgs.ChatMsg) {
	// Get requesting IP
	requestAddr := strings.Split(conn.RemoteAddr().String(), ":")
	requestIP := requestAddr[0]

	roomName, ok := ip2roomName[requestIP]
	if !ok {
		// TODO: handle error when ip is not in a room
	}

	r := rooms[roomName]

	fromIdx := uint64(binary.LittleEndian.Uint64(recvMsg.Payload))
	messages, nextIdx := r.FetchMessages(fromIdx)
	respIdx := make([]byte, 8)
	binary.LittleEndian.PutUint64(respIdx, nextIdx)
	payload := append(messages, respIdx...)

	response := msgs.ChatMsg{
		Status:  msgs.StatusOK,
		Payload: payload,
	}
	respond(conn, response)
	if fromIdx != nextIdx {
		username := ip2user[requestIP]
		log.Printf("%s (%s) got %s room's messages from %d up to message number %d\n", username, requestIP, roomName, fromIdx, nextIdx)
	}
}

func getMsgIdx(conn net.Conn, recvMsg msgs.ChatMsg) {
	// Get requesting IP
	requestAddr := strings.Split(conn.RemoteAddr().String(), ":")
	requestIP := requestAddr[0]

	roomName, ok := ip2roomName[requestIP]
	if !ok {
		// TODO: handle error when ip is not in a room
	}

	msgIdx := rooms[roomName].GetMsgIdx()

	payload := make([]byte, 8)
	binary.LittleEndian.PutUint64(payload, msgIdx)
	response := msgs.ChatMsg{
		Status:  msgs.StatusOK,
		Payload: payload,
	}
	respond(conn, response)
	// log.Printf("%s requested %s's msgIdx and got %d as response %v", ip2user[requestIP], roomName, msgIdx, payload)
}

func leave(conn net.Conn) {
	// Get requesting IP
	requestAddr := strings.Split(conn.RemoteAddr().String(), ":")
	requestIP := requestAddr[0]

	roomName, ok := ip2roomName[requestIP]
	if !ok {
		// TODO: handle error when ip is not in a room
	}

	r := rooms[roomName]

	r.LeaveUser(requestIP)

	response := msgs.ChatMsg{
		Status: msgs.StatusOK,
	}
	respond(conn, response)

	username := ip2user[requestIP]
	log.Printf("%s (%s) left %s room\n", username, requestIP, roomName)
}

func createRoom(conn net.Conn, recvMsg msgs.ChatMsg) {
	// Get requesting IP
	requestAddr := strings.Split(conn.RemoteAddr().String(), ":")
	requestIP := requestAddr[0]

	roomName := string(recvMsg.Payload)

	var response msgs.ChatMsg

	_, roomExists := rooms[roomName]
	if !roomExists {
		rooms[roomName] = &room{
			roomName:  roomName,
			userCount: 0,
		}
		rooms[roomName].users = make(map[string]struct{})
		response.Status = msgs.StatusOK
		username := ip2user[requestIP]
		log.Printf("%s (%s) created a new room named %s\n", username, requestIP, roomName)
	} else {
		response.Status = msgs.StatusRoomAlreadyExists
	}

	respond(conn, response)
}

func logout(conn net.Conn) {
	// Get requesting IP
	requestAddr := strings.Split(conn.RemoteAddr().String(), ":")
	requestIP := requestAddr[0]

	roomName, inRoom := ip2roomName[requestIP]
	if inRoom {
		r := rooms[roomName]
		r.LeaveUser(requestIP)
	}

	username := ip2user[requestIP]

	var response msgs.ChatMsg
	delete(ip2user, requestIP)
	delete(user2ip, username)

	response = msgs.ChatMsg{
		Status: msgs.StatusOK,
	}
	respond(conn, response)
	log.Printf("%s (%s) logged out", username, requestIP)
}

func respond(conn net.Conn, response msgs.ChatMsg) {
	respJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("error marshaling response: %s", err.Error())
	}
	conn.Write(respJSON)
}
