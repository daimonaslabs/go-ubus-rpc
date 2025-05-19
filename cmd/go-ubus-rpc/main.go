package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
)

func main() {
	ctx := context.TODO()
	opts := client.ClientOptions{Username: "root", Password: "D@!monas", URL: "http://10.0.0.1/ubus", Timeout: session.DefaultSessionTimeout}
	rpc, err := client.NewUbusRPC(ctx, &opts)

	if err != nil {
		fmt.Println(rpc, err)
	}
	uciGetOpts := client.UCIGetOptions{Config: "firewall", Section: "cfg04ad58"} //, Option: "src"}

	rpc.Call = rpc.UCI().Get(uciGetOpts)
	//sessionLoginOpts := client.SessionLoginOptions{Username: "root", Password: "D@!monas"}
	//rpc.Call = rpc.Session().Login(&sessionLoginOpts)
	//rpc.Call = rpc.UCI().Configs()
	response, err := rpc.Do(ctx)
	if err != nil {
		fmt.Println(err)
	}
	result, err := uciGetOpts.GetResult(response)
	//result, err := sessionLoginOpts.GetResult(response)
	//result, err := client.UCIConfigsOptions{}.GetResult(response)
	fmt.Println("response: ", response)
	fmt.Println("result: ", reflect.TypeOf(result), result)
	fmt.Println("err: ", reflect.TypeOf(err), err)
}
