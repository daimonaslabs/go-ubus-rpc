package client

import (
	"encoding/json"
	"errors"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
)

type SessionInterface interface {
	Login(opts *SessionLoginOptions) Call
}

// implements SessionInterface
// implements CallInterface
type sessionCall struct {
	Call
}

func newSessionCall(u *UbusRPC) *sessionCall {
	u.Call.setPath("session")
	return &sessionCall{u.Call}
}

func (c *sessionCall) Login(opts *SessionLoginOptions) Call {
	c.setProcedure("login")
	c.setSignature(opts)

	return c.Call
}

/*
################################################################
#
# all xOptions types are in this block. they all implement the
# Signature interface.
#
################################################################
*/

// implements Signature interface
type SessionLoginOptions struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  uint   `json:"timeout"`
}

func (SessionLoginOptions) isOptsType() {}

func (opts SessionLoginOptions) GetResult(p Response) (u LoginResult, err error) {
	if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case sessionResult:
			json.Unmarshal(data, &u)
		default:
			return LoginResult{}, errors.New("not a LoginResult")
		}
	} else { // error
		return LoginResult{}, errors.New(p[0].(ExitCode).Error())
	}
	return u, nil
}

/*
################################################################
#
# all xResult types are in this block.
#
################################################################
*/

// result of a `session login` command
type LoginResult struct {
	session.Session `json:",inline"`
}

// implements ResultObject interface
// used for handling the raw RPC response
type sessionResult struct {
	session.Session
}

func (sessionResult) isResultObject() {}

/*
################################################################
#
# all matchXResult funcs are in this block. used in init().
#
################################################################
*/

// checker for sessionResponse
func matchSessionResult(data json.RawMessage) (ResultObject, error) {
	var val session.Session

	if err := json.Unmarshal(data, &val); err == nil {
		if val.SessionID != "" { // easiest way to see if it unmarshaled into an empty Session struct
			return sessionResult{val}, nil
		}
	}
	return nil, nil
}
