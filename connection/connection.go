package connection

type Connection struct {
	Id             int
	Type           ConnectionType
	commandChan    chan<- string
	outputChan     <-chan string
	disconnectChan chan<- bool
}
