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
type Signature any

type UbusInterface interface {
	SessionCallGetter
	// UCICallGetter
}

// implements UbusInterface
type UbusRPC struct {
	Call CallInterface
	*clientset
	sessionCall
	//uciCall
}

func (u *UbusRPC) Session() SessionInterface {
	return newSessionCall(u)
}

type CallInterface interface {
	AsParams() Params
	SetSessionID(id SessionID)
	SetPath(p string)
	SetProcedure(p string)
	SetSignature(sig any)
}

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

type ResultInterface interface {
	//ToMap() map[string]any
}

// implements ResultInterface
type Result struct {
	resultJSON
}

type resultJSON = []resultContent

type resultContent struct {
	int
	resultContentInterface
}

type resultContentInterface interface {
}

//type SignatureInterface interface {
//	//SessionSignatureGetter
//	// UCISignatureGetter
//	ToMap() map[string]any
//}

//func (s *Signature) Signature() SignatureInterface {
//	return &LoginOptions{}
//}
//
//type SessionSignatureGetter interface {
//	Session() SessionSignatureInterface
//}
//
//type SessionSignatureInterface interface {
//	Login() map[string]any
//	Get() map[string]any
//}

//func (l *LoginOptions) Signature() SignatureInterface {
//	return Signature(map[string]any{
//		"username": l.Username,
//		"password": l.Password,
//		"timeout":  l.Timeout,
//	})
//}

//func ToMap[S SignatureInterface](s S) map[string]any {
//	print(s)
//	var sigMap map[string]any
//	data, err := json.Marshal(s)
//	if err != nil {
//		panic(err)
//	}
//
//	err = json.Unmarshal(data, &sigMap)
//	if err != nil {
//		panic(err)
//	}
//	return sigMap
//}

//func (s Signature) ToMap() map[string]any {
//	print(s)
//	var sigMap map[string]any
//	data, err := json.Marshal(s)
//	if err != nil {
//		panic(err)
//	}
//
//	err = json.Unmarshal(data, &sigMap)
//	if err != nil {
//		panic(err)
//	}
//	return sigMap
//}
