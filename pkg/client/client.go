package client

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

// add this later when we need a daemon client
//type ubusSessionManager struct {
//	Keepalive int
//}

// pass Keepalive to time.Ticker{}, goroutine listen to chan and refresh session
// when Ticker = Keepalive (timeout / 2 by default)
//func (u *UbusRPC) KeepAlive() {
//	go ...
//}

type clientset struct {
	rpcClient   *rpc.Client
	ubusSession *Session
}

type ClientOptions struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  uint   `json:"timeout"`
	URL      string `json:"url"`
}

func NewUbusRPC(ctx context.Context, opts *ClientOptions) (*UbusRPC, error) {
	c, err := newClientset(ctx, opts)
	return &UbusRPC{
		clientset: c,
	}, err
}

func (u *UbusRPC) Do(ctx context.Context) (r *Response, err error) {
	err = u.clientset.rpcClient.CallContext(ctx, &r, "call", u.Call.AsParams()...)
	return r, err
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
		ubusSession: &Session{},
	}
	login := &Call{
		SessionID: LoginSessionID,
		Path:      "session",
		Procedure: "login",
		Signature: map[string]any{
			"username": opts.Username,
			"password": opts.Password,
			"timeout":  opts.Timeout,
		},
	}

	// initialize ubus client
	result := Response{}
	request := login.AsParams()
	if err != nil {
		log.Fatalln(err)
	}
	err = c.rpcClient.CallContext(ctx, &result, "call", request...)
	if err != nil {
		log.Fatalln(err)
	}

	session, ok := result[1].(SessionResult)
	if !ok {
		err = errors.New("invalid response to ubus session login")
	}
	c.ubusSession = &session.Session
	return c, err
}
