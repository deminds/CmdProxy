package controller

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"

	"github.com/deminds/CmdProxy/model"
	"github.com/deminds/CmdProxy/session/types"
)

func (o *HttpController) ConsoleConnectHandler(respWriter http.ResponseWriter, request *http.Request) {
	logPrefix := "ConsoleConnectHandler()"
	glog.Info("%v Handle url: %v", logPrefix, request.URL.Path)

	if request.Method != http.MethodGet {
		glog.Errorf("%v Wrong message type. Expected: GET. Actual: %v", logPrefix, request.Method)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	sess, err := types.NewConsoleSession(o.idGenerator, o.timeoutSec)
	if err != nil {
		glog.Errorf("%v Error create local console connection. Error: %v", logPrefix, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	sess.Connect()

	o.sessionPool.Put(sess)

	response := model.ConnectResponse{
		SessionId: sess.GetId(),
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		glog.Errorf("%v Error marshal ConnectionResponse to json. sessID: %v Error: %v", logPrefix, sess.GetId(), err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write(responseBytes)
}

func (o *HttpController) ConsoleListHandler(respWriter http.ResponseWriter, request *http.Request) {
	logPrefix := "ConsoleListHandler()"
	glog.Info("%v Handle url: %v", logPrefix, request.URL.Path)
}
