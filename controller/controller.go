package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/deminds/CmdProxy/generatorid"
	"github.com/deminds/CmdProxy/model"
	"github.com/deminds/CmdProxy/session"
	"github.com/golang/glog"
)

const (
	SessionIdParam           = "sessionid"
	ContentTypeHeader        = "Content-Type"
	ContentTypeAppJsonHeader = "application/json"
)

func NewHttpController(
	pool *session.SessionPool,
	idGenerator *generatorid.IDGenerator,
	timeoutSec int) *HttpController {

	return &HttpController{
		sessionPool: pool,
		idGenerator: idGenerator,

		timeoutSec: timeoutSec,
	}
}

type HttpController struct {
	sessionPool *session.SessionPool
	idGenerator *generatorid.IDGenerator

	timeoutSec int
}

func (o *HttpController) DisconnectHandler(respWriter http.ResponseWriter, request *http.Request) {
	logPrefix := "DisconnectHandler()"
	glog.Infof("%v Handle url: %v", logPrefix, request.URL.Path)

	if request.Method != http.MethodGet {
		glog.Errorf("%v Wrong message type. Expected: GET. Actual: %v", logPrefix, request.Method)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	sessID := request.URL.Query().Get(SessionIdParam)
	if sessID == "" {
		glog.Errorf("%v Param %v not found in GET params", logPrefix, SessionIdParam)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := o.sessionPool.RemoveAndClose(sessID); err != nil {
		glog.Errorf("%v Error remove connection from sessionPool. "+
			"ID: %v, Error: %v", logPrefix, sessID, err)
		respWriter.WriteHeader(http.StatusNotFound)

		return
	}

	respWriter.WriteHeader(http.StatusOK)
}

func (o *HttpController) CommandHandler(respWriter http.ResponseWriter, request *http.Request) {
	logPrefix := "CommandHandler()"
	glog.Infof("%v Handle url: %v", logPrefix, request.URL.Path)

	if request.Method != http.MethodPost {
		glog.Errorf("%v Wrong message type. Expected: POST. Actual: %v", logPrefix, request.Method)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	contentTypeHeader := request.Header.Get(ContentTypeHeader)
	if contentTypeHeader != ContentTypeAppJsonHeader {
		glog.Errorf("%v Content-Type should be application/json. Content-Type: %v", logPrefix, contentTypeHeader)
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
		glog.Errorf("%v Error unmarshal to CommandRequest. RawMsg: %s, Error: %v", logPrefix, msgReqBytes, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	if msgReq.SessionId == "" || msgReq.Command == "" {
		glog.Errorf("%v Received not valid message. Msg: %+v", logPrefix, msgReq)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	sess, err := o.sessionPool.Get(msgReq.SessionId)
	if err != nil {
		glog.Errorf("%v SessionPool.Get(%v, %v). "+
			"Error: %v", logPrefix, session.SessionTypeConsole, msgReq.SessionId, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	cmdOutput, err := sess.Command(msgReq.Command)
	if err != nil {
		glog.Errorf("%v Error execute command. "+
			"ID: %v, Type: %v, CommandID: %v, Command: %v, Error: %v",
			logPrefix, sess.GetId(), sess.GetType(), msgReq.CommandId, msgReq.Command, err)
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
		glog.Errorf("%v Error marshal CommandResponse to json. RawMsg: %+v, Error: %v", logPrefix, msgResp, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	respWriter.WriteHeader(http.StatusOK)
	_, err = respWriter.Write(msgRespByte)
	if err != nil {
		glog.Errorf("%v Error write response POST message. Error: %v", logPrefix, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}
}
