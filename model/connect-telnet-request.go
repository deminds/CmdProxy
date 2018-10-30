package model

import "github.com/golang/glog"

type ConnectTelnetRequest struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Login    string `json:"login"`
	Password string `json:"password"`

	LoginExpectedString           string `json:"loginExpectedString"`
	PasswordExpectedString        string `json:"passwordExpectedString"`
	HostnameExpectedString        string `json:"hostnameExpectedString"`
	ContinueCommandExpectedString string `json:"continueCommandExpectedString"`
}

func (o *ConnectTelnetRequest) IsValid() bool {
	if o.Host == "" ||
		o.Port == 0 ||
		o.Login == "" ||
		o.Password == "" ||
		o.LoginExpectedString == "" ||
		o.PasswordExpectedString == "" ||
		o.HostnameExpectedString == "" {

		glog.Errorf("ConnectTelnetRequest.IsValid(). Is not valid. Struct: %+v", o)

		return false
	}

	return true
}
