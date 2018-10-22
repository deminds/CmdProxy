package model

type CommandRequest struct {
	SessionId string `json:"sessionId"`
	CommandId int    `json:"commandId"`
	Command   string `json:"command"`
}
