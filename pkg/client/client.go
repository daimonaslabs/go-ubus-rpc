package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client/ubus"
	"github.com/ethereum/go-ethereum/rpc"
)

type UbusSession struct {
	SessionID ubus.SessionID
	Keepalive int
}

type Client struct {
	RPCClient   *rpc.Client
	UbusSession *UbusSession
}

func (c *Client) NewUbusCall() *ubus.UbusCall {
	return &ubus.UbusCall{
		SessionID: c.UbusSession.SessionID,
	}
}

type ClientOptions struct {
	Username string
	Password string
	URL      string
	Timeout  int
}

// TODO Keepalive = time.Ticker{}, goroutine to poll and refresh session
func NewClient(ctx context.Context, opts *ClientOptions) (c *Client) {
	// initialize RPC client
	tokenHeader := rpc.WithHeader("Content-Type", "application/json")
	httpClient := rpc.WithHTTPClient(&http.Client{
		Timeout: 10 * time.Second,
	})
	rpcClient, err := rpc.DialOptions(ctx, opts.URL, httpClient, tokenHeader)
	if err != nil {
		log.Fatalln(err)
	}
	response := []any{}

	c = &Client{
		RPCClient: rpcClient,
		UbusSession: &UbusSession{
			SessionID: ubus.LoginSessionID,
			Keepalive: ubus.DefaultSessionTimeout / 2,
		},
	}
	login := &ubus.UbusCall{
		SessionID: c.UbusSession.SessionID,
		Path:      "session",
		Procedure: "login",
		Signature: ubus.Signature{
			"username": opts.Username,
			"password": opts.Password,
			"timeout":  opts.Timeout,
		},
	}
	request := login.ToParams()
	requestJSON, err := json.Marshal(request)
	if err != nil {
		log.Fatalln(err)
	}
	err = c.RPCClient.CallContext(ctx, &response, "call", request...)
	if err != nil {
		fmt.Println("err: ", err)
		log.Println("Error initiating ubus RPC session")
	}
	fmt.Println("request: ", request)
	fmt.Println("requestJSON: ", string(requestJSON))
	fmt.Println("response: ", response)

	return c
}
