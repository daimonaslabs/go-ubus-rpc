package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
)

func main() {
	ctx := context.TODO()
	opts := client.ClientOptions{Username: "root", Password: "D@!monas", URL: "http://10.0.0.1/ubus", Timeout: client.DefaultSessionTimeout}
	rpc, err := client.NewUbusRPC(ctx, &opts)

	if err != nil {
		fmt.Println(rpc, err)
	}

	//uciOpts := client.UCIOptions{Config: "firewall", Section: "cfg04ad58", Option: "src"}
	//rpc.Call = rpc.UCI().Get(&uciOpts)
	//rpc.Call = rpc.UCI().Configs()
	loginOpts := client.LoginOptions{Username: "root", Password: "D@!monas"}
	rpc.Call = rpc.Session().Login(&loginOpts)
	response, err := rpc.Do(ctx)
	fmt.Println(response, err)
	result, err := rpc.Session().GetResult(response)
	fmt.Println("result: ", reflect.TypeOf(result), result)
	fmt.Println("err: ", reflect.TypeOf(err), err)
}
