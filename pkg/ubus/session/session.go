package session

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus"
	"github.com/ethereum/go-ethereum/rpc"
)

type Data map[string]string

type Session struct {
	SessionID ubus.SessionID
	Timeout   int
	Expires   int
	//ACLs	ACL
	Data Data
}

type SessionLoginHandler struct {
	*ubus.UbusCall
	ubus.UbusResponseWriter
}

type SessionLoginResponse struct {
	Session
}

func (slr SessionLoginResponse) Write(resp []byte) error {
	resp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(resp, slr)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

type SessionLoginOptions struct {
	Timeout int
}

// make the actual call to ubus
// TODO is SessionLogin a special use case? should it be built into SessionManager?
func SessionLogin(uc *ubus.UbusCall, urw ubus.UbusResponseWriter) error {
	tokenHeader := rpc.WithHeader("Content-Type", "application/json")
	httpClient := rpc.WithHTTPClient(&http.Client{
		Timeout: 10 * time.Second,
	})
	ctx := context.Background()
	client, err := rpc.DialOptions(ctx, url, httpClient, tokenHeader)
	response := []any{}
	request := SessionCall{}.GetParams(Signature{"username": username, "password": password})
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
	return nil
}

func GetSessionLoginHandler(username, password string, opts *SessionLoginOptions) *SessionLoginHandler {
	if opts != nil {
		return &SessionLoginHandler{
			&ubus.UbusCall{
				SessionID: "00000000000000000000000000000000",
				Path:      "session",
				Procedure: "login",
				Signature: ubus.Signature{
					"username": username,
					"password": password,
					"timeout":  opts.Timeout,
				},
			},
			SessionLoginResponse{},
		}
	} else {
		return &SessionLoginHandler{
			&ubus.UbusCall{
				SessionID: "00000000000000000000000000000000",
				Path:      "session",
				Procedure: "login",
				Signature: ubus.Signature{
					"username": username,
					"password": password,
				},
			},
			SessionLoginResponse{},
		}
	}
}

type Params []any

func (c SessionCall) GetParams(sig Signature) Params {
	fmt.Println(sig)
	_, ok := sig["username"]
	if !ok {
		fmt.Println("Invalid signature, expecting: [ username: string, password: string ]")
	}
	_, ok = sig["password"]
	if !ok {
		fmt.Println("Invalid signature, expecting: [ username: string, password: string ]")
	}
	if len(sig) > 2 {

	}
	return Params{
		"00000000000000000000000000000000",
		"session", // Path
		"login",   // Procedure
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
	request := SessionCall{}.GetParams(Signature{"username": username, "password": password})
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
