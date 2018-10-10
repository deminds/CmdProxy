package connection

type TelnetConnection Connection

func (o *TelnetConnection) Connect(host string, port int, login string, password string) (*TelnetConnection, error) {

	return nil, nil
}

func (o *TelnetConnection) Disconnect() error {

	return nil
}

func (o *TelnetConnection) Command(command string) (string, error) {

	return "", nil
}

func (o *TelnetConnection) IsAlive() bool {

	return false
}
