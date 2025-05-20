package client

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
	"github.com/ethereum/go-ethereum/rpc"
)

type clientset struct {
	rpcClient   *rpc.Client
	ubusSession *session.Session
}

type UbusRPC struct {
	Call
	*clientset
	sessionRPC
	uciRPC
}

func (u *UbusRPC) Session() SessionInterface {
	return newSessionRPC(u)
}

func (u *UbusRPC) UCI() UCIInterface {
	return newUCIRPC(u)
}

func NewUbusRPC(ctx context.Context, opts *ClientOptions) (*UbusRPC, error) {
	c, err := newClientset(ctx, opts)
	return &UbusRPC{
		Call: Call{
			SessionID: c.ubusSession.SessionID,
		},
		clientset: c,
	}, err
}

func (u *UbusRPC) do(ctx context.Context) (r Response, err error) {
	err = u.clientset.rpcClient.CallContext(ctx, &r, "call", u.Call.asParams()...)
	return r, err
}

type ClientOptions struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  uint   `json:"timeout"`
	URL      string `json:"url"`
}

func newClientset(ctx context.Context, opts *ClientOptions) (c *clientset, err error) {
	// initialize RPC client
	tokenHeader := rpc.WithHeader("Content-Type", "application/json")
	httpClient := rpc.WithHTTPClient(&http.Client{
		Timeout: 10 * time.Second,
	})
	rpcClient, err := rpc.DialOptions(ctx, opts.URL, httpClient, tokenHeader)
	if err != nil {
		log.Fatalln(err)
	}

	c = &clientset{
		rpcClient:   rpcClient,
		ubusSession: &session.Session{},
	}
	loginOpts := SessionLoginOptions{
		Username: opts.Username,
		Password: opts.Password,
		Timeout:  opts.Timeout,
	}
	login := &Call{}
	login.setSessionID(session.LoginSessionID)
	login.setPath("session")
	login.setProcedure("login")
	login.setSignature(loginOpts)

	// initialize ubus client
	response := Response{}
	if err != nil {
		log.Fatalln(err)
	}

	err = c.rpcClient.CallContext(ctx, &response, "call", login.asParams()...)
	if err != nil {
		log.Fatalln(err)
	}

	session := response[1].(sessionResult)
	c.ubusSession = &session.Session
	return c, err
}
