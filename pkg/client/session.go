package client

import (
	"encoding/json"
)

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
	GetResult(r Response) SessionResult
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

func (c *sessionCall) GetResult(r Response) SessionResult {
	return r[1].(SessionResult)
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

// implements ResultObject interface
type SessionResult struct {
	Session
}

func (SessionResult) isResultObject() {}

// checker for SessionResponse
func matchSessionResult(data json.RawMessage) (ResultObject, error) {
	var val Session

	if err := json.Unmarshal(data, &val); err == nil {
		if val.SessionID != "" { // easiest way to see if it unmarshaled into an empty Session struct
			return SessionResult{val}, nil
		}
	}
	return nil, nil
}
