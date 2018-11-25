package main

import (
	"fmt"
	"time"
)

type message struct {
	msg    string
	sender string
	time   time.Time
}

type room struct {
	roomName  string
	userCount int
	messages  []message
}

func (r *room) JoinUser(ip string) {

}

func (r *room) String() string {
	ret := fmt.Sprintf("%s - %d users online", r.roomName, r.userCount)
	return ret

}
