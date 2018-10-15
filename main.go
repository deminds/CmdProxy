package main

import (
	"flag"
	"fmt"
	"github.com/deminds/CmdProxy/connection"
	"net/http"
	"os"
	"runtime/debug"

	httpHandlers "github.com/deminds/CmdProxy/provider"
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

	connectionPool := connection.NewConnectionPool()

	localHttpController := httpHandlers.LocalHttpController{
		ConnectionPool: connectionPool,
	}

	h := http.NewServeMux()

	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/connect", API_VERSION), httpHandlers.TelnetConnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/disconnect", API_VERSION), httpHandlers.TelnetDisconnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/command", API_VERSION), httpHandlers.TelnetCommandHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/telnet/list", API_VERSION), httpHandlers.TelnetListHandler)

	h.HandleFunc(fmt.Sprintf("/api/%v/local/connect", API_VERSION), localHttpController.LocalConnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/local/disconnect", API_VERSION), localHttpController.LocalDisconnectHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/local/command", API_VERSION), localHttpController.LocalCommandHandler)
	h.HandleFunc(fmt.Sprintf("/api/%v/local/list", API_VERSION), localHttpController.LocalListHandler)

	glog.Infof("Start listen %v:%v", HttpHost, HttpPort)
	l := http.ListenAndServe(fmt.Sprintf("%v:%v", HttpHost, HttpPort), h)

	glog.Fatal(l)
}
