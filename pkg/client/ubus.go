package client

import (
	"encoding/json"
)

const (
	LoginSessionID         SessionID = "00000000000000000000000000000000"
	DefaultSessionTimeout  uint      = 300
	NoExpirySessionTimeout uint      = 0
)

// maybe use this to do validation on the SessionID
//type sessionID [32]byte

type Params []any
type SessionID string

type UbusInterface interface {
	Session() SessionInterface
	UCI() UCIInterface
}

type Signature interface {
	isOptsType()
}

// implements UbusInterface
type UbusRPC struct {
	Call CallInterface
	*clientset
	sessionCall
	uciCall
}

func (u *UbusRPC) Session() SessionInterface {
	return newSessionCall(u)
}

func (u *UbusRPC) UCI() UCIInterface {
	return newUCICall(u)
}

type CallInterface interface {
	asParams() Params
	setSessionID(id SessionID)
	setPath(p string)
	setProcedure(p string)
	setSignature(sig Signature)
}

// implements CallInterface
// implements ResultTypeGetter
type Call struct {
	SessionID SessionID
	Path      string
	Procedure string
	Signature Signature
}

func (uc *Call) asParams() Params {
	return Params{uc.SessionID, uc.Path, uc.Procedure, uc.Signature}
}

func (uc *Call) setSessionID(id SessionID) {
	uc.SessionID = id
}

func (uc *Call) setPath(p string) {
	uc.Path = p
}

func (uc *Call) setProcedure(p string) {
	uc.Procedure = p
}

func (uc *Call) setSignature(sig Signature) {
	data, err := json.Marshal(sig)
	if err != nil {
		panic(err)
	}
	uc.Signature = sig
	json.Unmarshal(data, &uc.Signature)
}
