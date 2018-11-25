package msgs

import (
	"encoding/json"
)

//ChatMsg is the data format used for communication between client and server
type ChatMsg struct {
	Command int
	Status  int
	Payload []byte
}

//Command code definitions
const (
	CmdLogin int = iota + 1
	CmdList
	CmdJoin
	CmdMsg
	CmdFetch
	CmdGetMsgIdx
	CmdLeave
	CmdCreateRoom
	CmdLogout
)

//Response Status code definitions
const (
	StatusOK int = iota + 1
	StatusInexistentRoom
	StatusUserNotInRoom
	StatusRoomAlreadyExists
)

func (msg ChatMsg) String() string {
	prettyJSON, _ := json.MarshalIndent(msg, "", "    ")
	return string(prettyJSON) + "\n"
}
