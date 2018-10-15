package connection

type Connection struct {
	id             int
	connectionType ConnectionType
	commandChan    chan<- string
	outputChan     <-chan string
	disconnectChan chan<- bool
}

type IConnection interface {
	Disconnect() error
	Command(command string) (string, error)
	IsAlive() bool
	GetId() int
	SetId(id int)
	GetType() ConnectionType
	IsClose() bool
}
