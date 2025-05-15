package client

import (
	"encoding/json"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
)

type Params []any

type UbusInterface interface {
	Session() SessionInterface
	UCI() UCIInterface
}

// all implementations have an implicit method of GetResult(Response) (xResult, error)
type Signature interface {
	isOptsType()
}

type CallInterface interface {
	AsParams() Params
	setSessionID(id session.SessionID)
	setPath(p string)
	setProcedure(p string)
	setSignature(sig Signature)
}

// implements CallInterface
// implements ResultTypeGetter
type Call struct {
	SessionID session.SessionID
	Path      string
	Procedure string
	Signature Signature
}

func (uc *Call) asParams() Params {
	return Params{uc.SessionID, uc.Path, uc.Procedure, uc.Signature}
}

func (uc *Call) setSessionID(id session.SessionID) {
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
