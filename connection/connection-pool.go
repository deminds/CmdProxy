package connection

type ConnectionPool struct {
	connections map[ConnectionType]map[int]*Connection
}

func (o *ConnectionPool) Get(connectionId int, connectionType ConnectionType) *Connection {
	return nil
}

func (o *ConnectionPool) TelnetConnect(host string, port int, login string, password string) (int, error) {
	return 0, nil
}

func (o *ConnectionPool) LocalConnect() (int, error) {
	return 0, nil
}
