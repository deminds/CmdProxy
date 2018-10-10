package model

type CommandResponse struct {
	CommandRequest `json:",inline"`
	Output         string `json:"output"`
}
