package client

import (
	"encoding/json"
)

var LoginSessionID = SessionID([]byte("00000000000000000000000000000000"))
var DefaultSessionTimeout = int(300)

type SessionID [32]byte
type Params []any

type Signature any // Signature struct types must have JSON tags
type UbusResponse map[string]any

type UbusInterface interface {
	SessionCallGetter
	// UCICallGetter
}

type SessionCallGetter interface {
	Session() SessionInterface
}

// implements UbusInterface
type UbusRPC struct {
	*clientset
	sessionCall
	//uciCall
}

func (u *UbusRPC) Session() SessionInterface {
	return newSessionCall(u)
}

func newSessionCall(u *UbusRPC) *sessionCall {
	u.sessionCall.SetSessionID(u.ubusSession.SessionID)
	u.sessionCall.SetPath("session")
	return &u.sessionCall
}

type SessionInterface interface {
	Login(opts *LoginOptions) *sessionCall
}

// implements SessionInterface
type sessionCall struct {
	Call
}

func (c *sessionCall) Login(opts *LoginOptions) *sessionCall {
	c.SetProcedure("login")
	c.SetSignature(map[string]any{
		"username": opts.Username,
		"password": opts.Password,
		"timeout":  opts.Timeout,
	})

	return c
}

type LoginOptions struct {
	Username string
	Password string
	Timeout  int
}

//type CallInterface interface {
//	AsParams() Params
//	SetSessionID(id SessionID)
//	SetPath(p string)
//	SetProcedure(p string)
//	SetSignature(sig any)
//}

// implements CallInterface
type Call struct {
	SessionID SessionID
	Path      string
	Procedure string
	Signature Signature
}

func (uc *Call) SetSessionID(id SessionID) {
	uc.SessionID = id
}

func (uc *Call) SetPath(p string) {
	uc.Path = p
}

func (uc *Call) SetProcedure(p string) {
	uc.Procedure = p
}

func (uc *Call) SetSignature(sig any) {
	data, err := json.Marshal(sig)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &uc.Signature)
}

func (uc *Call) AsParams() Params {
	return Params{uc.SessionID, uc.Path, uc.Procedure, uc.Signature}
}

//type Session struct {
//	SessionID string
//	Timeout   int
//	Expires   int
//	//ACLs	ACL
//	Data map[string]string
//}
