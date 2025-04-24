package session

import (
	"github.com/daimonaslabs/go-ubus-rpc/pkg/client/ubus"
)

type SessionInterface interface {
	Login(opts *LoginOptions) (*ubus.UbusResponse, error)
}

type sessionCall struct {
	*ubus.UbusCall
}

func loginCall(opts *LoginOptions) sessionCall {
	call := sessionCall{}
	call.SetPath("session")
	call.SetProcedure("login")
	call.SetSignature(opts)

	return call
}

type LoginOptions struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  int    `json:"timeout,omitempty"`
}

/*
type SessionLoginHandler struct {
	*ubus.UbusCall
	ubus.UbusResponseWriter
}

type SessionLoginResponse struct {
	Session
}

func (slr SessionLoginResponse) Write(resp []byte) error {
	resp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(resp, slr)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}
*/
