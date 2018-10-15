package provider

import (
	"encoding/json"
	"github.com/deminds/CmdProxy/connection"
	"github.com/deminds/CmdProxy/model"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	ConnectionIdParam = "connectionid"
)

func NewLocalHttpController(connPool *connection.ConnectionPool) *LocalHttpController {
	return &LocalHttpController{
		ConnectionPool: connPool,
	}
}

type LocalHttpController struct {
	ConnectionPool *connection.ConnectionPool
}

func (o *LocalHttpController) LocalConnectHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("Http handle local/connect")

	if request.Method != http.MethodGet {
		glog.Errorf("Wrong message type. Expected: GET. Actual: %v", request.Method)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	conn, err := connection.ConnectLocal()
	if err != nil {
		glog.Errorf("Error create local connection. Error: %v", err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	o.ConnectionPool.Put(conn)

	response := model.ConnectionResponse{
		ConnectionId: conn.GetId(),
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		glog.Errorf("Error marshal ConnectionResponse to json. Error: %v", err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write(responseBytes)
}

func (o *LocalHttpController) LocalDisconnectHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("Http handle local/disconnect")

	if request.Method != http.MethodGet {
		glog.Errorf("Wrong message type. Expected: GET. Actual: %v", request.Method)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	connectionId, err := strconv.Atoi(request.URL.Query().Get(ConnectionIdParam))
	if err != nil {
		glog.Errorf("Error param '%v' in GET request. Error: %v", ConnectionIdParam, err)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}
	glog.Infoln("after parse")

	if err := o.ConnectionPool.Remove(connection.ConnectionTypeLocal, connectionId); err != nil {
		glog.Errorf("Error remove connection from connectionPool. "+
			"ConnectionId: %v, Error: %v", connectionId, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	glog.Infoln("after remove")

	respWriter.WriteHeader(http.StatusOK)
}

func (o *LocalHttpController) LocalStatusHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("Http handle local/status")

}

func (o *LocalHttpController) LocalCommandHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("Http handle local/command")

	if request.Method != http.MethodPost {
		glog.Errorf("Wrong message type. Expected: POST. Actual: %v", request.Method)
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	msgReqBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		glog.Errorf("Error read POST message. Error: %v", err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	var msgReq model.CommandRequest
	if err := json.Unmarshal(msgReqBytes, &msgReq); err != nil {
		glog.Errorf("Error unmarshal to CommandRequest. RawMsg: %s, Error: %v", msgReqBytes, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	conn, err := o.ConnectionPool.Get(connection.ConnectionTypeLocal, msgReq.ConnectionId)
	if err != nil {
		glog.Errorf("ConnectionPool.Get(%v, %v). "+
			"Error: %v", connection.ConnectionTypeLocal, msgReq.ConnectionId, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}
	cmdOutput, err := conn.Command(msgReq.Command)

	msgResp := model.CommandResponse{
		Output: cmdOutput,
		CommandRequest: model.CommandRequest{
			Command:      msgReq.Command,
			CommandId:    msgReq.CommandId,
			ConnectionId: msgReq.ConnectionId,
		},
	}

	msgRespByte, err := json.Marshal(msgResp)
	if err != nil {
		glog.Errorf("Error marshal CommandResponse to json. RawMsg: %+v, Error: %v", msgResp, err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	respWriter.WriteHeader(http.StatusOK)
	_, err = respWriter.Write(msgRespByte)
	if err != nil {
		glog.Errorf("Error write response POST message. Error: %v", err)
		respWriter.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (o *LocalHttpController) LocalListHandler(respWriter http.ResponseWriter, request *http.Request) {
	glog.Info("Http handle local/list")

}
