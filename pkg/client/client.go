package client

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/rpc"
)

type ubusSession struct {
	SessionID SessionID
	Keepalive int
}

type clientset struct {
	rpcClient   *rpc.Client
	ubusSession *ubusSession
}

type ClientOptions struct {
	Username string
	Password string
	URL      string
	Timeout  int
}

//type Caller interface {
//	MakeCall() (*ubus.UbusResponse, error)
//}
//
//func (rpc *UbusRPC) MakeCall() (resp *ubus.UbusResponse, err error) {
//	err = rpc.client.rpcClient.CallContext(context.TODO(), resp, "call", rpc.call.AsParams()...)
//	return &ubus.UbusResponse{}, err
//}
//
//func NewCaller(opts *ClientOptions) Caller {
//	return &UbusRPC{
//		client: newClient(context.TODO(), opts),
//		call:   ubus.UbusCall{},
//	}
//}

func NewUbusRPC(opts *ClientOptions) (*UbusRPC, error) {
	c, err := newClientset(opts)
	return &UbusRPC{
		clientset: c,
	}, err
}

// TODO add opts for the call to the context somehow
func (u *UbusRPC) Do(ctx context.Context) (r *UbusResponse, err error) {
	err = u.clientset.rpcClient.Call(&r, "call", u.Session().Login(&LoginOptions{Username: "ha", Password: "haha"}))
	return &UbusResponse{}, err
}

// TODO Keepalive = time.Ticker{}, goroutine to poll and refresh session
func newClientset(opts *ClientOptions) (c *clientset, err error) {
	// initialize RPC client
	//tokenHeader := rpc.WithHeader("Content-Type", "application/json")
	//httpClient := rpc.WithHTTPClient(&http.Client{
	//	Timeout: 10 * time.Second,
	//})
	//rpcClient, err := rpc.DialOptions(ctx, opts.URL, httpClient, tokenHeader)
	rpcClient, err := rpc.Dial(opts.URL)
	if err != nil {
		log.Fatalln(err)
	}
	response := []any{}

	c = &clientset{
		rpcClient: rpcClient,
		ubusSession: &ubusSession{
			SessionID: LoginSessionID,
			Keepalive: opts.Timeout / 2,
		},
	}
	login := &Call{
		SessionID: c.ubusSession.SessionID,
		Path:      "session",
		Procedure: "login",
		Signature: map[string]any{
			"username": opts.Username,
			"password": opts.Password,
			"timeout":  opts.Timeout,
		},
	}
	request := login.AsParams()
	if err != nil {
		log.Fatalln(err)
	}
	err = c.rpcClient.Call(&response, "call", request...)
	if err != nil {
		log.Fatalln(err)
	}

	return c, err
}
