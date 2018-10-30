package types

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/deminds/CmdProxy/generatorid"
	"github.com/deminds/CmdProxy/model"
	"github.com/deminds/CmdProxy/session"
	"github.com/golang/glog"
	"github.com/ziutek/telnet"
)

const (
	ContinueCommand = " "
)

func NewTelnetSession(idGenerator *generatorid.IDGenerator, timeoutSec int, requestData model.ConnectTelnetRequest) (*TelnetSession, error) {
	id, err := idGenerator.Next()
	if err != nil {
		return nil, fmt.Errorf("NewTelnetSession(). Generate id. Error: %v", err)
	}

	sess := &TelnetSession{
		id:          id,
		isClose:     false,
		sessionType: session.SessionTypeTelnet,
		timeout:     time.Duration(timeoutSec) * time.Second,

		loginExpectedString:    requestData.LoginExpectedString,
		passwordExpectedString: requestData.PasswordExpectedString,
		hostnameExpectedString: requestData.HostnameExpectedString,
		continueExpectedString: requestData.ContinueCommandExpectedString,

		host: requestData.Host,
		port: requestData.Port,

		login:    requestData.Login,
		password: requestData.Password,

		output:     make(chan string),
		command:    make(chan string),
		disconnect: make(chan bool),
	}

	glog.Infof("NewConsoleSession() Host: %v, Port: %v, ID: %v, Type: %v, Timeout: %v",
		sess.host, sess.port, sess.id, sess.sessionType, sess.timeout)

	return sess, nil
}

// TODO: remove channels
type TelnetSession struct {
	id          string
	isClose     bool
	sessionType session.SessionType
	timeout     time.Duration

	loginExpectedString    string
	passwordExpectedString string
	hostnameExpectedString string
	continueExpectedString string

	host string
	port int

	login    string
	password string

	sess *telnet.Conn

	command    chan string
	output     chan string
	disconnect chan bool
}

func (o *TelnetSession) Connect() error {
	logPrefix := "TelnetSession.Connect()"

	addr := fmt.Sprintf("%v:%v", o.host, o.port)

	glog.Infof("%v Addr: %v, ID: %v, Type: %v", logPrefix, addr, o.id, o.sessionType)

	sess, err := telnet.Dial("tcp", addr)
	if err != nil {
		o.isClose = true

		return fmt.Errorf("%v Error telnet.Dial(). ID: %v, Type: %v, Addr: %v, Error: %v",
			logPrefix, o.id, o.sessionType, addr, err)
	}

	sess.SetUnixWriteMode(true)
	o.sess = sess

	// login
	resp, err := o.readUntil(o.loginExpectedString)
	if err != nil {
		return fmt.Errorf("%v Read after connect. Wait: %v, Error: %v", logPrefix, o.loginExpectedString, err)
	}
	glog.Infof("%v Connect. Response: %v", logPrefix, resp)

	if err := o.sendLine(o.login); err != nil {
		return fmt.Errorf("%v Send login. Error: %v", logPrefix, err)
	}

	resp, err = o.readUntil(o.passwordExpectedString)
	if err != nil {
		return fmt.Errorf("%v Read after send login. Error: %v", logPrefix, err)
	}
	glog.Infof("%v Send login. Response: %v", logPrefix, resp)

	if err := o.sendLine(o.password); err != nil {
		return fmt.Errorf("%v Send password. Error: %v", logPrefix, err)
	}

	resp, err = o.readUntil(o.hostnameExpectedString)
	if err != nil {
		return fmt.Errorf("%v Read after send password. Wait: %v, Error: %v", logPrefix, o.hostnameExpectedString, err)
	}
	glog.Infof("%v Send password. Response: %v", logPrefix, resp)

	go o.start()

	return nil
}

func (o *TelnetSession) Command(command string) (string, error) {
	logPrefix := "TelnetSession.Command()"

	glog.Infof("%v Execute command. "+
		"ID: %v, Type: %v, Command: '%v'", logPrefix, o.id, o.sessionType, command)

	if o.isClose {
		return "", fmt.Errorf("%v Session is close. "+
			"ID: %v, Type: %v, Command: %v", logPrefix, o.id, o.sessionType, command)
	}

	o.command <- command

	select {
	case res := <-o.output:
		glog.Infof("%v Received output. "+
			"ID: %v, Type: %v, Output: '%s'", logPrefix, o.id, o.sessionType, res)

		return res, nil
	case <-time.After(o.timeout):
		o.isClose = true

		return "", fmt.Errorf("%v Timeout wait output. "+
			"ID: %v, Type: %v", logPrefix, o.id, o.sessionType)
	}
}

func (o *TelnetSession) Ping() bool {
	if _, err := o.Command(PingCommand); err != nil {
		return false
	}

	return true
}

func (o *TelnetSession) GetId() string {
	return o.id
}

func (o *TelnetSession) GetType() session.SessionType {
	return o.sessionType
}

func (o *TelnetSession) IsClose() bool {
	return o.isClose
}

func (o *TelnetSession) Close() {
	glog.Infof("TelnetSession.Close(). ID: %v, Type: %v", o.id, o.sessionType)

	o.isClose = true
	o.disconnect <- true
}

func (o *TelnetSession) start() {
	logPrefix := "TelnetSession.start()"
	for {
		select {
		case cmd, ok := <-o.command:
			if !ok {
				glog.Infof("%v Command chan was close. Exit routine. "+
					"ID: %v, Type: %v", logPrefix, o.id, o.sessionType)
				o.isClose = true

				return
			}

			cmd = strings.Trim(cmd, " ")
			if cmd == "" {
				o.output <- EmptyCommandMsg

				continue
			}

			if err := o.sendLine(cmd); err != nil {
				glog.Errorf("%v Send command. Exit routine. Command: %v, Error: %v", logPrefix, cmd, err)
				o.isClose = true

				return
			}

			resp, err := o.readStringUntil(o.hostnameExpectedString)
			if err != nil {
				glog.Errorf("%v Read after send command. Exit routine. Wait: %v, Error: %v", logPrefix, o.hostnameExpectedString, err)
				o.isClose = true

				return
			}

			if !o.isClose {
				o.output <- resp
			} else {
				glog.Info("%v Session was closed. Exit routine. Id: %v, Type: %v", logPrefix, o.id, o.sessionType)
				o.isClose = true

				return
			}

		case d := <-o.disconnect:
			if d {
				glog.Infof("%v Received disconnect request. "+
					"ID: %v, Type: %v", logPrefix, o.id, o.sessionType)
				o.isClose = true

				return
			}

		case <-time.After(o.timeout):
			glog.Infof("%v Timeout between commands was reach. Drop session. "+
				"ID: %v, Type: %v", logPrefix, o.id, o.sessionType)
			o.isClose = true

			return
		}
	}
}

// Will find delim in full output
func (o *TelnetSession) readUntil(delim string) (string, error) {
	logPrefix := "TelnetSession.readUntil()"

	o.sess.SetReadDeadline(time.Now().Add(o.timeout))
	resBytes, err := o.sess.ReadUntil(delim)
	if err != nil {
		return "", fmt.Errorf("%v o.sess.ReadUntil() "+
			"ID: %v, Delim: %v, Error: %v", logPrefix, o.id, delim, err)
	}

	return string(resBytes), nil
}

// Will find delim in begin of string
func (o *TelnetSession) readStringUntil(delim string) (string, error) {
	logPrefix := "TelnetSession.readStringUntil()"

	// append next line '\n' before delimiter
	delim = string([]byte{10}) + delim

	delims := []string{}
	if o.continueExpectedString != "" {
		delims = append(delims, o.continueExpectedString)
	}
	delims = append(delims, delim)

	buf := bytes.Buffer{}
	for {
		o.sess.SetReadDeadline(time.Now().Add(o.timeout))
		resBytes, idx, err := o.sess.ReadUntilIndex(delims...)
		if err != nil {
			return "", fmt.Errorf("%v o.sess.ReadUntilIndex() "+
				"ID: %v, Delim: %v, Error: %v", logPrefix, o.id, delim, err)
		}

		if _, err := buf.Write(resBytes); err != nil {
			return "", fmt.Errorf("%v buf.Write() "+
				"ID: %v, Delim: %v, Error: %v", logPrefix, o.id, delim, err)
		}

		if o.continueExpectedString == "" || idx != 0 {
			break
		}

		if err := o.sendLine(" "); err != nil {
			return "", fmt.Errorf("%v sendLine() "+
				"ID: %v, Delim: %v, Error: %v", logPrefix, o.id, delim, err)
		}
	}

	return buf.String(), nil
}

func (o *TelnetSession) sendLine(cmd string) error {
	buf := make([]byte, len(cmd)+1)

	copy(buf, cmd)
	buf[len(buf)-1] = '\n'

	o.sess.SetWriteDeadline(time.Now().Add(o.timeout))

	_, err := o.sess.Write(buf)
	if err != nil {
		return fmt.Errorf("TelnetSession.sendLine() o.sess.Write(). "+
			"ID: %v, Command: '%v', Error: %v", o.id, cmd, err)
	}

	return nil
}
