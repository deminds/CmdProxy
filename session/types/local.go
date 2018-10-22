package types

import (
	"fmt"
	"github.com/deminds/CmdProxy/generatorid"
	"github.com/deminds/CmdProxy/session"
	"github.com/golang/glog"
	"time"
)

const (
	TimeoutSec  = 10
	PingCommand = "\n"
)

func NewLocalSession(idGenerator *generatorid.IDGenerator) (*LocalSession, error) {
	id, err := idGenerator.Next()
	if err != nil {
		return nil, fmt.Errorf("NewLocalSession(). Generate id. Error: %v", err)
	}

	sess := &LocalSession{
		id:          id,
		isClose:     false,
		sessionType: session.SessionTypeLocal,

		output:     make(chan string),
		command:    make(chan string),
		disconnect: make(chan bool),
	}

	return sess, nil
}

type LocalSession struct {
	id          string
	isClose     bool
	sessionType session.SessionType

	command    chan string
	output     chan string
	disconnect chan bool
}

func (o *LocalSession) Connect() {
	glog.Infof("LocalSession.Connect() ID: %v, Type: %v", o.id, o.sessionType)
	go o.start()
}

func (o *LocalSession) Command(command string) (string, error) {
	glog.Infof("LocalSession.Command(%v). Execute command. "+
		"ID: %v, Type: %v", command, o.id, o.sessionType)

	if o.isClose {
		return "", fmt.Errorf("LocalSession.Command(%v). Session is close. "+
			"ID: %v, Type: %v", command, o.id, o.sessionType)
	}

	glog.Info("Command() Before execute command")
	o.command <- command
	glog.Info("Command() After execute command")
	select {
	case res := <-o.output:
		glog.Infof("LocalSession.Command(%v). Received output. "+
			"ID: %v, Type: %v, Output: %v", command, o.id, o.sessionType, res)

		return res, nil
	case <-time.After((TimeoutSec + 10) * time.Second):
		glog.Infof("Command() Case timeout. Set isClose=true. return")
		o.isClose = true
		return "", fmt.Errorf("LocalSession.Command(%v) Timeout wait output. "+
			"ID: %v, Type: %v", command, o.id, o.sessionType)
	}
}

func (o *LocalSession) Ping() bool {
	if _, err := o.Command(PingCommand); err != nil {
		return false
	}

	return true
}

func (o *LocalSession) GetId() string {
	return o.id
}

func (o *LocalSession) GetType() session.SessionType {
	return o.sessionType
}

func (o *LocalSession) IsClose() bool {
	return o.isClose
}

func (o *LocalSession) Close() {
	glog.Infof("LocalSession.Close(). ID: %v, Type: %v", o.id, o.sessionType)

	o.disconnect <- true
	o.isClose = true
}

func (o *LocalSession) start() {
	for {
		select {
		case c, ok := <-o.command:
			if !ok {
				glog.Infof("LocalSession.start() Command chan was close. "+
					"ID: %v, Type: %v", o.id, o.sessionType)
				o.isClose = true

				return
			}

			glog.Info("start() case o.command. Before sleep")
			time.Sleep(15 * time.Second) // test sleep
			glog.Info("start() case o.command. After sleep")
			if !o.isClose {
				o.output <- fmt.Sprintf("Echo command:\n  %s", c)
			} else {
				glog.Info("start() Session was closed. return")
				return
			}

		case d := <-o.disconnect:
			glog.Info("start() case o.disconnect")
			if d {
				glog.Infof("LocalSession.start() Received disconnect request. "+
					"ID: %v, Type: %v", o.id, o.sessionType)
				o.isClose = true

				return
			}

		case <-time.After((TimeoutSec - 1) * time.Second):
			glog.Info("start() case timeout (-1). Set isClose=true. return")
			glog.Infof("LocalSession.start() Timeout between commands was reach. Drop session. "+
				"ID: %v, Type: %v", o.id, o.sessionType)
			o.isClose = true

			return
		}
	}
}
