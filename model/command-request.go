package model

type CommandRequest struct {
	ConnectionId int    `json:"connectionId"`
	CommandId    int    `json:"commandId"`
	Command      string `json:"command"`
}
