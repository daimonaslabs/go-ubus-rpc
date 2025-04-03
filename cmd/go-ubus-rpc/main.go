package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

type Params []any

//type InitSessionRequest struct {
//	JSONRPC string `json:"jsonrpc"`
//	ID      int    `json:"id"`
//	Method  string `json:"method"`
//	Params  []any  `json:"params"`
//}

func NewInitSessionRequest(username, password string) Params {
	//	return &InitSessionRequest{
	//		JSONRPC: "2.0",
	//		ID:      1,
	//		Method:  "call",
	//		Params: []any{
	//			"00000000000000000000000000000000",
	//			"session",
	//			"login",
	//			map[string]string{
	//				"username": username,
	//				"password": password,
	//			},
	//		},
	//	}
	return []any{
		"00000000000000000000000000000000",
		"session",
		"login",
		map[string]string{
			"username": username,
			"password": password,
		},
	}
}

func main() {
	tokenHeader := rpc.WithHeader("Content-Type", "application/json")
	httpClient := rpc.WithHTTPClient(&http.Client{
		Timeout: 10 * time.Second,
	})
	ctx := context.Background()
	client, err := rpc.DialOptions(ctx, "http://10.0.0.1/ubus", httpClient, tokenHeader)
	response := []any{}
	request := NewInitSessionRequest("root", "D@!monas")
	requestJSON, err := json.Marshal(request)
	err = client.CallContext(ctx, &response, "call", request...)
	if err != nil {
		fmt.Println("err: ", err)
		log.Println("Error initiating ubus RPC session")
	}
	fmt.Println("request: ", request)
	fmt.Println("requestJSON: ", string(requestJSON))
	fmt.Println("response: ", response)

}
