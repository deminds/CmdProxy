package main

import (
	"flag"
	"fmt"
	"github.com/deminds/CmdProxy/generatorid"
	"github.com/deminds/CmdProxy/session"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/deminds/CmdProxy/controller"
	"github.com/golang/glog"
)

const (
	API_VERSION = "v1.0"
)

var (
	HttpHost = flag.String("host", "0.0.0.0", "IP for start application on it")
	HttpPort = flag.Int("port", 25505, "Port for start application on it")

	consoleTimeoutSec = flag.Int("consoleTimeout", 10, "Set timeout for console session and timeout for command in session")
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			glog.Errorf("main()\n%+v", r, debug.Stack())
		}
		glog.Infoln(">>>>> Service stop")
		glog.Flush()
	}()
	flag.Parse()

	glog.Infof(">>>>> Service start. Args: %+v", os.Args)

	idGenerator := generatorid.NewIDGenerator()

	pool := session.NewSessionPool()

	h := http.NewServeMux()

	consoleHttpController := controller.NewConsoleHttpController(pool, idGenerator, *consoleTimeoutSec)

	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/connect", API_VERSION), controller.TelnetConnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/disconnect", API_VERSION), controller.TelnetDisconnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/command", API_VERSION), controller.TelnetCommandHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/list", API_VERSION), controller.TelnetListHandler)

	h.HandleFunc(fmt.Sprintf("/api/%v/console/connect", API_VERSION), consoleHttpController.ConsoleConnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/console/disconnect", API_VERSION), consoleHttpController.ConsoleDisconnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/console/command", API_VERSION), consoleHttpController.ConsoleCommandHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/console/list", API_VERSION), consoleHttpController.ConsoleListHandler)

	glog.Infof("Start listen %v:%v", *HttpHost, *HttpPort)
	l := http.ListenAndServe(fmt.Sprintf("%v:%v", *HttpHost, *HttpPort), h)

	glog.Fatal(l)
}
