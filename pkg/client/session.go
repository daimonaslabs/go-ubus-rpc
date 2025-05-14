package client

import (
	"encoding/json"
	"errors"
)

const SessionResultType = "session"

type Session struct {
	SessionID SessionID `json:"ubus_rpc_session"`
	Timeout   int       `json:"timeout"`
	Expires   int       `json:"expires"`
	ACLs      ACL       `json:"acls"`
	Data      Data      `json:"data"`
}

type ACL struct {
	AccessGroup map[string][]string `json:"access-group"`
	CGIIO       map[string][]string `json:"cgi-io,omitempty"`
	File        map[string][]string `json:"file,omitempty"`
	Ubus        map[string][]string `json:"ubus"`
	UCI         map[string][]string `json:"uci,omitempty"`
}

type Data struct {
	Username string `json:"username"`
}

type SessionInterface interface {
	GetResult(p Response) (u SessionResult, err error)
	Login(opts *LoginOptions) CallInterface
}

// implements SessionInterface
// implements CallInterface
type sessionCall struct {
	Call
}

func newSessionCall(u *UbusRPC) *sessionCall {
	u.sessionCall.setSessionID(u.ubusSession.SessionID)
	u.sessionCall.setPath("session")
	return &u.sessionCall
}

func (c *sessionCall) GetResult(p Response) (u SessionResult, err error) {
	if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case sessionResult:
			u.Type = SessionResultType
			json.Unmarshal(data, &u.Values)
			return u, nil
		default:
			return SessionResult{}, errors.New("not a SessionResult")
		}
	} else { // error
		return SessionResult{}, errors.New(p[0].(ExitCode).Error())
	}
}

func (c *sessionCall) Login(opts *LoginOptions) CallInterface {
	c.setProcedure("login")
	c.setSignature(opts)

	return c
}

// implements Signature interface
type LoginOptions struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  uint   `json:"timeout"`
}

func (LoginOptions) isOptsType() {}

type SessionResult struct {
	resultStatic
}

// implements ResultObject interface
type sessionResult struct {
	Session
}

func (sessionResult) isResultObject() {}

// checker for SessionResponse
func matchSessionResult(data json.RawMessage) (ResultObject, error) {
	var val Session

	if err := json.Unmarshal(data, &val); err == nil {
		if val.SessionID != "" { // easiest way to see if it unmarshaled into an empty Session struct
			return sessionResult{val}, nil
		}
	}
	return nil, nil
}
