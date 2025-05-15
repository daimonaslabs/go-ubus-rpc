package session

const (
	LoginSessionID         SessionID = "00000000000000000000000000000000"
	DefaultSessionTimeout  uint      = 300
	NoExpirySessionTimeout uint      = 0
)

// maybe use this to do validation on the SessionID
//type sessionID [32]byte

type SessionID string

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
