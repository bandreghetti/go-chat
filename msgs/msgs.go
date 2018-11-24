package msgs

import (
	"encoding/json"
)

//Command code definitions
const (
	CmdLogin int = iota + 1
)

//Response Status code definitions
const (
	StatusOK int = iota + 1
)

type ChatMsg struct {
	Command int
	Status  int
	Payload []byte
}

func (msg ChatMsg) String() string {
	prettyJSON, _ := json.MarshalIndent(msg, "", "    ")
	return string(prettyJSON) + "\n"
}
