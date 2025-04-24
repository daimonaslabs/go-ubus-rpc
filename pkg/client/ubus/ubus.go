package ubus

import (
	"encoding/json"
)

var LoginSessionID = SessionID([]byte("00000000000000000000000000000000"))
var DefaultSessionTimeout = int(300)

type SessionID [32]byte
type Params []any

type Signature any // Signature struct types must have JSON tags
type UbusResponse map[string]any

type UbusCall struct {
	SessionID SessionID
	Path      string
	Procedure string
	Signature any
}

func (uc *UbusCall) SetSessionID(id SessionID) {
	uc.SessionID = id
}

func (uc *UbusCall) SetPath(p string) {
	uc.Path = p
}

func (uc *UbusCall) SetProcedure(p string) {
	uc.Procedure = p
}

func (uc *UbusCall) SetSignature(sig any) {
	data, err := json.Marshal(sig)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &uc.Signature)
}

func (uc *UbusCall) AsParams() Params {
	return Params{uc.SessionID, uc.Path, uc.Procedure, uc.Signature}
}

type Session struct {
	SessionID string
	Timeout   int
	Expires   int
	//ACLs	ACL
	Data map[string]string
}
