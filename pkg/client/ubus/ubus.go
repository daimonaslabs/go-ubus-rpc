package ubus

var LoginSessionID = SessionID("00000000000000000000000000000000")
var DefaultSessionTimeout = int(300)

type SessionID string
type Path string
type Procedure string
type Signature map[string]any
type Params []any

type Data map[string]string

type UbusCall struct {
	SessionID SessionID
	Path      Path
	Procedure Procedure
	Signature Signature
}

func (uc *UbusCall) ToParams() Params {
	return Params{uc.SessionID, uc.Path, uc.Procedure, uc.Signature}
}

type Session struct {
	SessionID string
	Timeout   int
	Expires   int
	//ACLs	ACL
	Data Data
}
