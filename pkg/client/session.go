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

func newSessionCall(u *UbusRPC) *sessionCall {
	u.sessionCall.SetSessionID(u.ubusSession.SessionID)
	u.sessionCall.SetPath("session")
	return &u.sessionCall
}

type SessionInterface interface {
	Login(opts *LoginOptions) CallInterface
}

// implements SessionInterface
type sessionCall struct {
	Call
}

func (c *sessionCall) Login(opts *LoginOptions) CallInterface {
	c.SetProcedure("login")
	c.SetSignature(opts)

	return c
}

type LoginOptions struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  uint   `json:"timeout"`
}

// implements ResultObject interface
type SessionResult struct {
	Session
}

func (SessionResult) isResultObject() {}

// checker for SessionResponse
func matchSessionResult(data json.RawMessage) (ResultObject, error) {
	var val Session

	if err := json.Unmarshal(data, &val); err == nil {
		if val.SessionID != "" { // easiest way to see if it unmarshaled into a empty Session struct
			return SessionResult{val}, nil
		}
	}
	return nil, nil
}

//		if reflect.TypeOf(val.SessionID).Kind() == reflect.Map {
//			fmt.Println("not a sessionResult!")
//			return nil, nil
//		}
//	}
