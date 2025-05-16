package client

import (
	"encoding/json"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
)

type Params []any

//type UbusInterface interface {
//	Session() SessionInterface
//	UCI() UCIInterface
//}

// all implementations have an implicit method of GetResult(Response) (xResult, error)
type Signature interface {
	isOptsType()
}

//type CallInterface interface {
//	asParams() Params
//	setSessionID(id session.SessionID)
//	setPath(p string)
//	setProcedure(p string)
//	setSignature(sig Signature)
//}

// implements CallInterface
// implements ResultTypeGetter
type Call struct {
	SessionID session.SessionID
	Path      string
	Procedure string
	Signature Signature
}

func (c *Call) asParams() Params {
	return Params{c.SessionID, c.Path, c.Procedure, c.Signature}
}

func (c *Call) setSessionID(id session.SessionID) {
	c.SessionID = id
}

func (c *Call) setPath(p string) {
	c.Path = p
}

func (c *Call) setProcedure(p string) {
	c.Procedure = p
}

func (uc *Call) setSignature(sig Signature) {
	data, err := json.Marshal(sig)
	if err != nil {
		panic(err)
	}
	uc.Signature = sig
	json.Unmarshal(data, &uc.Signature)
}
