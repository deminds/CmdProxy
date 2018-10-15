package connection

import (
	"fmt"
	"github.com/deminds/CmdProxy/session"
	"time"
)

const (
	CommandTimeoutSec = 5
	IsAliveCommand    = "\n"
)

type LocalConnection struct {
	Connection
}

func ConnectLocal() (*LocalConnection, error) {
	command := make(chan string)
	output := make(chan string)
	disconnect := make(chan bool)

	connection := &LocalConnection{
		Connection{
			connectionType: ConnectionTypeLocal,
			commandChan:    command,
			outputChan:     output,
			disconnectChan: disconnect,
		},
	}

	go session.LocalSession(command, output, disconnect)

	return connection, nil
}

func (o *LocalConnection) Disconnect() error {
	o.disconnectChan <- true

	return nil
}

func (o *LocalConnection) Command(command string) (string, error) {
	o.commandChan <- command

	select {
	case result := <-o.outputChan:
		return result, nil
	case <-time.After(CommandTimeoutSec * time.Second):
		return "", fmt.Errorf("Timeout reached for wait response from session. Disconnect. "+
			"Id: %v, Type: %v, Command: %v", o.id, o.connectionType, command)
	}
}

func (o *LocalConnection) IsAlive() bool {
	_, err := o.Command(IsAliveCommand)
	if err != nil {
		return false
	}

	return true
}

func (o *LocalConnection) GetId() int {
	return o.id
}

func (o *LocalConnection) SetId(id int) {
	o.id = id
}

func (o *LocalConnection) GetType() ConnectionType {
	return o.connectionType
}

func (o *LocalConnection) IsClose() bool {
	if _, ok := <-o.outputChan; ok {
		return false
	}

	return true
}
