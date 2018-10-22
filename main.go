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
	HttpHost = *flag.String("host", "0.0.0.0", "IP for start application on it")
	HttpPort = *flag.Int("port", 25505, "Port for start application on it")
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

	glog.Infof(">>>>> Serivice start. Args: %+v", os.Args)

	idGenerator := generatorid.NewIDGenerator()

	pool := session.NewSessionPool()

	h := http.NewServeMux()

	localHttpController := controller.NewLocalHttpController(pool, idGenerator)

	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/connect", API_VERSION), controller.TelnetConnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/disconnect", API_VERSION), controller.TelnetDisconnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/command", API_VERSION), controller.TelnetCommandHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/list", API_VERSION), controller.TelnetListHandler)

	h.HandleFunc(fmt.Sprintf("/api/%v/local/connect", API_VERSION), localHttpController.LocalConnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/local/disconnect", API_VERSION), localHttpController.LocalDisconnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/local/command", API_VERSION), localHttpController.LocalCommandHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/local/list", API_VERSION), localHttpController.LocalListHandler)

	glog.Infof("Start listen %v:%v", HttpHost, HttpPort)
	l := http.ListenAndServe(fmt.Sprintf("%v:%v", HttpHost, HttpPort), h)

	glog.Fatal(l)
}
