package connection

type LocalConnection Connection

func (o *LocalConnection) Connect() (*TelnetConnection, error) {

	return nil, nil
}

func (o *LocalConnection) Disconnect() error {

	return nil
}

func (o *LocalConnection) Command(command string) (string, error) {

	return "", nil
}

func (o *LocalConnection) IsAlive() bool {

	return false
}
