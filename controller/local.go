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

func NewLocalHttpController(pool *session.SessionPool, idGenerator *generatorid.IDGenerator) *LocalHttpController {
	return &LocalHttpController{
		sessionPool: pool,
		idGenerator: idGenerator,
	}
}

type LocalHttpController struct {
	sessionPool *session.SessionPool
	idGenerator *generatorid.IDGenerator
}

func (o *LocalHttpController) LocalConnectHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("LocalConnectHandler()")

	if request.Method != http.MethodGet {
		glog.Errorf("LocalConnectHandler() Wrong message type. Expected: GET. Actual: %v", request.Method)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	sess, err := types.NewLocalSession(o.idGenerator)
	if err != nil {
		glog.Errorf("LocalConnectHandler() Error create local connection. Error: %v", err)
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
		glog.Errorf("LocalConnectHandler() Error marshal ConnectionResponse to json. sessID: %v Error: %v", sess.GetId(), err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write(responseBytes)
}

func (o *LocalHttpController) LocalDisconnectHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("LocalDisconnectHandler()")

	if request.Method != http.MethodGet {
		glog.Errorf("LocalDisconnectHandler() Wrong message type. Expected: GET. Actual: %v", request.Method)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	sessID := request.URL.Query().Get(SessionIdParam)
	if sessID == "" {
		glog.Errorf("LocalDisconnectHander() Param %v not found in GET params", SessionIdParam)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := o.sessionPool.RemoveAndClose(session.SessionTypeLocal, sessID); err != nil {
		glog.Errorf("LocalDisconnectHandler() Error remove connection from sessionPool. "+
			"ID: %v, Error: %v", sessID, err)
		respWriter.WriteHeader(http.StatusNotFound)

		return
	}

	respWriter.WriteHeader(http.StatusOK)
}

func (o *LocalHttpController) LocalStatusHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("Http handle local/status")

}

func (o *LocalHttpController) LocalCommandHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("LocalCommandHandler()")

	if request.Method != http.MethodPost {
		glog.Errorf("LocalCommandHandler() Wrong message type. Expected: POST. Actual: %v", request.Method)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	msgReqBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		glog.Errorf("LocalCommandHandler() Error read POST message. Error: %v", err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	var msgReq model.CommandRequest
	if err := json.Unmarshal(msgReqBytes, &msgReq); err != nil {
		glog.Errorf("LocalCommandHandler() Error unmarshal to CommandRequest. RawMsg: %s, Error: %v", msgReqBytes, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	if msgReq.SessionId == "" || msgReq.Command == "" {
		glog.Errorf("LocalCommandHandler() Received not valid message. Msg: %+v", msgReq)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	sess, err := o.sessionPool.Get(session.SessionTypeLocal, msgReq.SessionId)
	if err != nil {
		glog.Errorf("LocalCommandHandler() SessionPool.Get(%v, %v). "+
			"Error: %v", session.SessionTypeLocal, msgReq.SessionId, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}
	cmdOutput, err := sess.Command(msgReq.Command)

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
		glog.Errorf("LocalCommandHandler() Error marshal CommandResponse to json. RawMsg: %+v, Error: %v", msgResp, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	respWriter.WriteHeader(http.StatusOK)
	_, err = respWriter.Write(msgRespByte)
	if err != nil {
		glog.Errorf("LocalCommandHandler() Error write response POST message. Error: %v", err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (o *LocalHttpController) LocalListHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("Http handle local/list")

}
