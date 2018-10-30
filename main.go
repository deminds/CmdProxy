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

	sessionTimeoutSec = flag.Int("timeout", 10, "Set timeout for session and timeout for command in session")
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

	httpController := controller.NewHttpController(pool, idGenerator, *sessionTimeoutSec)

	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/connect", API_VERSION), httpController.TelnetConnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/list", API_VERSION), httpController.TelnetListHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/disconnect", API_VERSION), httpController.DisconnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/command", API_VERSION), httpController.CommandHandler)

	h.HandleFunc(fmt.Sprintf("/api/%v/console/connect", API_VERSION), httpController.ConsoleConnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/console/list", API_VERSION), httpController.ConsoleListHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/console/disconnect", API_VERSION), httpController.DisconnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/console/command", API_VERSION), httpController.CommandHandler)

	glog.Infof("Start listen %v:%v", *HttpHost, *HttpPort)
	l := http.ListenAndServe(fmt.Sprintf("%v:%v", *HttpHost, *HttpPort), h)

	glog.Fatal(l)
}
