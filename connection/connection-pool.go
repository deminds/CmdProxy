package connection

import (
	"fmt"
	"github.com/golang/glog"
	"sync"
)

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		connections:      map[ConnectionType]map[int]IConnection{},
		lastConnectionId: 0,
	}
}

type ConnectionPool struct {
	connections      map[ConnectionType]map[int]IConnection
	lastConnectionId int
	mutex            sync.Mutex
}

func (o *ConnectionPool) Get(connectionType ConnectionType, connectionId int) (IConnection, error) {
	typeToConnections, exist := o.connections[connectionType]
	if !exist {
		msg := fmt.Sprintf("Try get connection from connectionPool. "+
			"Connection not found. Id: %v, Type: %v", connectionId, connectionType)
		glog.Errorf(msg)

		return nil, fmt.Errorf(msg)
	}

	connection, exist := typeToConnections[connectionId]
	if !exist {
		msg := fmt.Sprintf("Try get connection from connectionPool. "+
			"Connection not found. Id: %v, Type: %v", connectionId, connectionType)
		glog.Errorf(msg)

		return nil, fmt.Errorf(msg)
	}

	if connection.IsClose() {
		msg := fmt.Sprintf("Get closed connection from connectionPool. Remove connection. "+
			"Id: %v, Type: %v", connection.GetId(), connection.GetType())
		glog.Errorf(msg)
		o.Remove(connectionType, connectionId)

		return nil, fmt.Errorf(msg)
	}

	return connection, nil
}

func (o *ConnectionPool) Put(connection IConnection) (int, error) {
	if connection.IsClose() {
		msg := fmt.Sprintf("Try put in connectionPool closed connection")
		glog.Errorf(msg)

		return 0, fmt.Errorf(msg)
	}

	connectionType := connection.GetType()
	if _, exist := o.connections[connectionType]; !exist {
		o.connections[connectionType] = map[int]IConnection{}
	}

	connectionId := o.getNextConnectionId()

	o.connections[connectionType][connectionId] = connection
	connection.SetId(connectionId)

	return connectionId, nil
}

func (o *ConnectionPool) Remove(connectionType ConnectionType, connectionId int) error {
	typeToConnections, exist := o.connections[connectionType]
	if !exist {
		return fmt.Errorf("Remove from connectionPool. "+
			"Type not found in connections. Type: %v, ConnectionId: %v", connectionType, connectionId)
	}

	connection, exist := typeToConnections[connectionId]
	if !exist {
		return fmt.Errorf("Remove from connectionPool. "+
			"ConnectionId not found in connections. Type: %v, ConnectionId: %v", connectionType, connectionId)
	}

	if !connection.IsClose() {
		connection.Disconnect()
	}

	delete(o.connections[connectionType], connectionId)

	return nil
}

func (o *ConnectionPool) TelnetConnect(host string, port int, login string, password string) (int, error) {
	return 0, nil
}

func (o *ConnectionPool) LocalConnect() (int, error) {
	return 0, nil
}

func (o *ConnectionPool) getNextConnectionId() int {
	o.mutex.Lock()
	c := o.lastConnectionId
	c++
	o.lastConnectionId = c
	o.mutex.Unlock()

	return c
}
