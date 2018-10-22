package session

type ISession interface {
	Connect()
	Command(command string) (string, error)
	Ping() bool
	GetId() string
	GetType() SessionType
	IsClose() bool
	Close()
}
