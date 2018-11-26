package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bandreghetti/go-chat/msgs"
)

type message struct {
	msg    string
	sender string
	time   time.Time
}

type room struct {
	roomName  string
	userCount uint
	messages  []message
	users     map[string]struct{}
}

func (r *room) Empty() bool {
	if r.userCount == 0 {
		return true
	}
	return false
}

func (r *room) JoinUser(ip string) {
	// TODO: check if user has already joined the room
	r.userCount++
	r.users[ip] = struct{}{}
	ip2roomName[ip] = r.roomName
	msgText := fmt.Sprintf("%s has joined the room.", ip2user[ip])
	message := message{
		msg:    msgText,
		sender: "server",
		time:   time.Now(),
	}
	r.messages = append(r.messages, message)
}

func (r *room) LeaveUser(ip string) {
	// TODO: check if user really is in the room
	r.userCount--
	delete(r.users, ip)
	delete(ip2roomName, ip)
	msgText := fmt.Sprintf("%s has left the room.", ip2user[ip])
	message := message{
		msg:    msgText,
		sender: "server",
		time:   time.Now(),
	}
	r.messages = append(r.messages, message)
}

func (r *room) ListUsers() []byte {
	var users []string
	if r.userCount > 0 {
		users = append(users, fmt.Sprintf("List of users in %s", r.roomName))
		for user := range r.users {
			users = append(users, ip2user[user])
		}
	} else {
		users = append(users, fmt.Sprintf("There are no users in %s", r.roomName))
	}
	return []byte(strings.Join(users, "\n"))
}

func (r *room) PostMessage(ip string, msgText string) int {
	if _, ok := r.users[ip]; !ok {
		// TODO: handle error when user is not in the room
		return msgs.StatusUserNotInRoom
	}
	message := message{
		msg:    msgText,
		sender: ip2user[ip],
		time:   time.Now(),
	}
	r.messages = append(r.messages, message)

	return msgs.StatusOK
}

func (r *room) FetchMessages(lastIdx uint64) ([]byte, uint64) {
	var messageSlice []string
	for i := lastIdx; i < uint64(len(r.messages)); i++ {
		msg := r.messages[i]
		_, month, day := msg.time.Date()
		messageString := fmt.Sprintf("[%02d/%02d %02d:%02d] %s: %s",
			day, month, msg.time.Hour(), msg.time.Minute(), msg.sender, msg.msg)
		messageSlice = append(messageSlice, messageString)
	}
	return []byte(strings.Join(messageSlice, "\n")), uint64(len(r.messages))
}

func (r *room) GetMsgIdx() uint64 {
	return uint64(len(r.messages))
}

func (r *room) String() string {
	ret := fmt.Sprintf("%s - %d users online", r.roomName, r.userCount)
	return ret
}
