package model

type CommandRequest struct {
	SessionId string `json:"sessionid"`
	CommandId int    `json:"commandid,omitempty"`
	Command   string `json:"command"`
}
