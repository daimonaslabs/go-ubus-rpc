package ubus

import (
	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/client/ubus"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
)

type UbusInterface interface {
	UbusRPC() client.Caller
	SessionSetter
}

type SessionSetter interface {
	Session()
}

type SessionInterface interface {
	Login(opts *session.LoginOptions) (*ubus.UbusResponse, error)
}

type UbusRPC struct {
	rpc *client.UbusRPC
}

func (uc *UbusRPC) Session() {
	uc.rpc.Call.Path = "session"
}

type sessionCall struct {
	ubus.UbusCall
}

func (c sessionCall) Login() {
	c.Procedure = "login"
}

//type sessionCall struct {
//	ubus.UbusCall
//}
//
//func (call sessionCall) sessionCall() {
//	call.Path = "session"
//}

func (call *ubusCall) LoginCall(opts *session.LoginOptions) (*ubus.UbusResponse, error) {
	call.SetPath("session")
	call.SetProcedure("login")
	call.SetSignature(opts)

	call.Call()
}

//type UCICallGetter interface {
//	UCI() *UbusCall
//}
