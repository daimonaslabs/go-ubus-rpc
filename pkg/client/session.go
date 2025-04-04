package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

type SessionID string
type Path string
type Procedure string
type Signature map[string]any
type Data map[string]string

type Params []any

type Session struct {
	SessionID SessionID
	Timeout   int
	Expires   int
	//ACLs	ACL
	Data Data
}

type UbusCall struct {
	SessionID SessionID
	Path      Path
	Procedure Procedure
	Signature Signature
}

type Caller interface {
	GetParams(sig Signature) Params
}

type SessionLoginCall struct {
	UbusCall
}

func (c SessionLoginCall) GetParams(sig Signature) Params {
	fmt.Println(sig)
	_, ok := sig["username"]
	if !ok {
		fmt.Println("Invalid signature, expecting: [ username: string, password: string ]")
	}
	_, ok = sig["password"]
	if !ok {
		fmt.Println("Invalid signature, expecting: [ username: string, password: string ]")
	}
	return Params{
		"00000000000000000000000000000000",
		"session",
		"login",
		sig,
	}
}

//func NewInitSessionRequest(username, password string) Params {
//	return []any{
//		"00000000000000000000000000000000",
//		"session",
//		"login",
//		map[string]string{
//			"username": username,
//			"password": password,
//		},
//	}
//}

func NewSession(username, password, url string) any {
	tokenHeader := rpc.WithHeader("Content-Type", "application/json")
	httpClient := rpc.WithHTTPClient(&http.Client{
		Timeout: 10 * time.Second,
	})
	ctx := context.Background()
	client, err := rpc.DialOptions(ctx, url, httpClient, tokenHeader)
	response := []any{}
	request := SessionLoginCall{}.GetParams(Signature{"username": username, "password": password})
	requestJSON, err := json.Marshal(request)
	err = client.CallContext(ctx, &response, "call", request...)
	if err != nil {
		fmt.Println("err: ", err)
		log.Println("Error initiating ubus RPC session")
	}
	fmt.Println("request: ", request)
	fmt.Println("requestJSON: ", string(requestJSON))
	fmt.Println("response: ", response)

	return response
}
