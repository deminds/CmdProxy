package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/deminds/CmdProxy/model"
	"github.com/deminds/CmdProxy/session/types"
	"github.com/golang/glog"
)

func (o *HttpController) TelnetConnectHandler(respWriter http.ResponseWriter, request *http.Request) {
	logPrefix := "TelnetConnectHandler()"
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
		glog.Errorf("%v Error read POST message. Error: %v", logPrefix, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	var msgReq model.ConnectTelnetRequest
	if err := json.Unmarshal(msgReqBytes, &msgReq); err != nil {
		glog.Errorf("%v Error unmarshal to ConnectTelnetRequest. RawMsg: %s, Error: %v", logPrefix, msgReqBytes, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}
	glog.Infof("%v Received POST: %+v", logPrefix, msgReq)

	if !msgReq.IsValid() {
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	sess, err := types.NewTelnetSession(o.idGenerator, o.timeoutSec, msgReq)
	if err != nil {
		glog.Errorf("%v NewTelnetSession(). Error: %v", logPrefix, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	if err := sess.Connect(); err != nil {
		glog.Errorf("%v sess.Connect() Error: %v", logPrefix, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	o.sessionPool.Put(sess)

	response := model.ConnectResponse{
		SessionId: sess.GetId(),
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		glog.Errorf("%v Error marshal ConnectionResponse to json. sessID: %v Error: %v", sess.GetId(), err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write(responseBytes)
}

func (o *HttpController) TelnetListHandler(respWriter http.ResponseWriter, request *http.Request) {
	logPrefix := "TelnetListHandler()"
	glog.Info("%v Handle url: %v", logPrefix, request.URL.Path)
}
