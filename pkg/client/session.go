package client

type Session struct {
	SessionID SessionID         `json:"ubus_rpc_session"`
	Timeout   int               `json:"timeout"`
	Expires   int               `json:"expires"`
	ACLs      ACL               `json:"acls"`
	Data      map[string]string `json:"data"`
}

type ACL struct {
	AccessGroup map[string][]string `json:"access-group"`
	CGIIO       map[string][]string `json:"cgi-io,omitempty"`
	File        map[string][]string `json:"file,omitempty"`
	Ubus        map[string][]string `json:"ubus"`
	UCI         map[string][]string `json:"uci,omitempty"`
}

type SessionCallGetter interface {
	Session() SessionInterface
}

func newSessionCall(u *UbusRPC) *sessionCall {
	u.sessionCall.SetSessionID(u.SessionID)
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
	c.SetSignature(map[string]any{
		"username": opts.Username,
		"password": opts.Password,
		"timeout":  opts.Timeout,
	})

	return c
}

type LoginOptions struct {
	Username string
	Password string
	Timeout  uint
}
