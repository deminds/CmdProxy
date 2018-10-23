package types

import (
	"fmt"
	"github.com/deminds/CmdProxy/generatorid"
	"github.com/deminds/CmdProxy/session"
	"github.com/golang/glog"
	"os/exec"
	"strings"
	"time"
)

const (
	PingCommand          = "\n"
	EmptyCommandMsg      = "command is empty string"
	CommandArgsSeparator = " "
)

func NewConsoleSession(idGenerator *generatorid.IDGenerator, timeoutSec int) (*ConsoleSession, error) {
	id, err := idGenerator.Next()
	if err != nil {
		return nil, fmt.Errorf("NewConsoleSession(). Generate id. Error: %v", err)
	}

	sess := &ConsoleSession{
		id:          id,
		isClose:     false,
		sessionType: session.SessionTypeConsole,
		timeout:     timeoutSec,

		output:     make(chan string),
		command:    make(chan string),
		disconnect: make(chan bool),
	}

	glog.Infof("NewConsoleSession() ID: %v, Type: %v, Timeout: %v", sess.id, sess.sessionType, sess.timeout)

	return sess, nil
}

type ConsoleSession struct {
	id          string
	isClose     bool
	sessionType session.SessionType
	timeout     int

	command    chan string
	output     chan string
	disconnect chan bool
}

func (o *ConsoleSession) Connect() {
	glog.Infof("ConsoleSession.Connect() ID: %v, Type: %v", o.id, o.sessionType)
	go o.start()
}

func (o *ConsoleSession) Command(command string) (string, error) {
	glog.Infof("ConsoleSession.Command(%v). Execute command. "+
		"ID: %v, Type: %v", command, o.id, o.sessionType)

	if o.isClose {
		return "", fmt.Errorf("ConsoleSession.Command(%v). Session is close. "+
			"ID: %v, Type: %v", command, o.id, o.sessionType)
	}

	o.command <- command

	select {
	case res := <-o.output:
		glog.Infof("ConsoleSession.Command(%v). Received output. "+
			"ID: %v, Type: %v, Output: %v", command, o.id, o.sessionType, res)

		return res, nil
	case <-time.After(time.Duration(o.timeout) * time.Second):
		o.isClose = true

		return "", fmt.Errorf("ConsoleSession.Command(%v) Timeout wait output. "+
			"ID: %v, Type: %v", command, o.id, o.sessionType)
	}
}

func (o *ConsoleSession) Ping() bool {
	if _, err := o.Command(PingCommand); err != nil {
		return false
	}

	return true
}

func (o *ConsoleSession) GetId() string {
	return o.id
}

func (o *ConsoleSession) GetType() session.SessionType {
	return o.sessionType
}

func (o *ConsoleSession) IsClose() bool {
	return o.isClose
}

func (o *ConsoleSession) Close() {
	glog.Infof("ConsoleSession.Close(). ID: %v, Type: %v", o.id, o.sessionType)

	o.disconnect <- true
	o.isClose = true
}

func (o *ConsoleSession) start() {
	for {
		select {
		case c, ok := <-o.command:
			if !ok {
				glog.Infof("ConsoleSession.start() Command chan was close. "+
					"ID: %v, Type: %v", o.id, o.sessionType)
				o.isClose = true

				return
			}

			c = strings.Trim(c, CommandArgsSeparator)
			if c == "" {
				o.output <- EmptyCommandMsg

				continue
			}

			if !o.isClose {
				cPaths := strings.Split(c, CommandArgsSeparator)

				cName := cPaths[0]
				cArgs := []string{}
				if len(cPaths) > 1 {
					cArgs = cPaths[1:]
				}

				glog.Infof("exec.Command(%v, %v)", cName, cArgs)

				cmd := exec.Command(cName, cArgs...)

				out, err := cmd.CombinedOutput()
				outStr := string(out)
				if err != nil {
					// TODO: put to output struct (out, err)
					glog.Errorf("start() cmd.Run() failed. ID: %v, Type: %v, Error: %v", o.id, o.sessionType, err)
					outStr = err.Error()
				}

				o.output <- outStr
			} else {
				glog.Info("start() Session was closed. Exit routine. Id: %v, Type: %v", o.id, o.sessionType)

				return
			}

		case d := <-o.disconnect:
			if d {
				glog.Infof("ConsoleSession.start() Received disconnect request. "+
					"ID: %v, Type: %v", o.id, o.sessionType)
				o.isClose = true

				return
			}

		case <-time.After(time.Duration(o.timeout) * time.Second):
			glog.Infof("ConsoleSession.start() Timeout between commands was reach. Drop session. "+
				"ID: %v, Type: %v", o.id, o.sessionType)
			o.isClose = true

			return
		}
	}
}
