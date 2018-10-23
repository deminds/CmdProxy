package controller

import (
	"encoding/json"
	"github.com/deminds/CmdProxy/generatorid"
	"github.com/deminds/CmdProxy/model"
	"github.com/deminds/CmdProxy/session"
	"github.com/deminds/CmdProxy/session/types"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
)

const (
	SessionIdParam = "sessionid"
)

func NewConsoleHttpController(
	pool *session.SessionPool,
	idGenerator *generatorid.IDGenerator,
	timeoutSec int) *ConsoleHttpController {

	return &ConsoleHttpController{
		sessionPool: pool,
		idGenerator: idGenerator,
		timeoutSec:  timeoutSec,
	}
}

type ConsoleHttpController struct {
	sessionPool *session.SessionPool
	idGenerator *generatorid.IDGenerator

	timeoutSec int
}

func (o *ConsoleHttpController) ConsoleConnectHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("ConsoleConnectHandler()")

	if request.Method != http.MethodGet {
		glog.Errorf("ConsoleConnectHandler() Wrong message type. Expected: GET. Actual: %v", request.Method)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	sess, err := types.NewConsoleSession(o.idGenerator, o.timeoutSec)
	if err != nil {
		glog.Errorf("ConsoleConnectHandler() Error create local console connection. Error: %v", err)
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
		glog.Errorf("ConsoleConnectHandler() Error marshal ConnectionResponse to json. sessID: %v Error: %v", sess.GetId(), err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write(responseBytes)
}

func (o *ConsoleHttpController) ConsoleDisconnectHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("ConsoleDisconnectHandler()")

	if request.Method != http.MethodGet {
		glog.Errorf("ConsoleDisconnectHandler() Wrong message type. Expected: GET. Actual: %v", request.Method)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	sessID := request.URL.Query().Get(SessionIdParam)
	if sessID == "" {
		glog.Errorf("ConsoleDisconnectHandler() Param %v not found in GET params", SessionIdParam)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := o.sessionPool.RemoveAndClose(session.SessionTypeConsole, sessID); err != nil {
		glog.Errorf("ConsoleDisconnectHandler() Error remove connection from sessionPool. "+
			"ID: %v, Error: %v", sessID, err)
		respWriter.WriteHeader(http.StatusNotFound)

		return
	}

	respWriter.WriteHeader(http.StatusOK)
}

func (o *ConsoleHttpController) ConsoleStatusHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("Http handle console/status")

}

func (o *ConsoleHttpController) ConsoleCommandHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("ConsoleCommandHandler()")

	if request.Method != http.MethodPost {
		glog.Errorf("ConsoleCommandHandler() Wrong message type. Expected: POST. Actual: %v", request.Method)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	msgReqBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		glog.Errorf("ConsoleCommandHandler() Error read POST message. Error: %v", err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	var msgReq model.CommandRequest
	if err := json.Unmarshal(msgReqBytes, &msgReq); err != nil {
		glog.Errorf("ConsoleCommandHandler() Error unmarshal to CommandRequest. RawMsg: %s, Error: %v", msgReqBytes, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	if msgReq.SessionId == "" || msgReq.Command == "" {
		glog.Errorf("ConsoleCommandHandler() Received not valid message. Msg: %+v", msgReq)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	sess, err := o.sessionPool.Get(session.SessionTypeConsole, msgReq.SessionId)
	if err != nil {
		glog.Errorf("ConsoleCommandHandler() SessionPool.Get(%v, %v). "+
			"Error: %v", session.SessionTypeConsole, msgReq.SessionId, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	cmdOutput, err := sess.Command(msgReq.Command)
	if err != nil {
		glog.Errorf("ConsoleCommandHandler() Error execute command. "+
			"ID: %v, Type: %v, CommandID: %v, Command: %v, Error: %v",
			sess.GetId(), sess.GetType(), msgReq.CommandId, msgReq.Command, err)
		respWriter.WriteHeader(http.StatusNotModified)

		return
	}

	msgResp := model.CommandResponse{
		Output: cmdOutput,
		CommandRequest: model.CommandRequest{
			Command:   msgReq.Command,
			CommandId: msgReq.CommandId,
			SessionId: sess.GetId(),
		},
	}

	msgRespByte, err := json.Marshal(msgResp)
	if err != nil {
		glog.Errorf("ConsoleCommandHandler() Error marshal CommandResponse to json. RawMsg: %+v, Error: %v", msgResp, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	respWriter.WriteHeader(http.StatusOK)
	_, err = respWriter.Write(msgRespByte)
	if err != nil {
		glog.Errorf("ConsoleCommandHandler() Error write response POST message. Error: %v", err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (o *ConsoleHttpController) ConsoleListHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("Http handle console/list")

}
