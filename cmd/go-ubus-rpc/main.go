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

	uciOpts := client.UCIOptions{Config: "firewall", Section: "cfg04ad58", Option: "src"}
	rpc.Call = rpc.UCI().Get(&uciOpts)
	response, err := rpc.Do(ctx)
	fmt.Println(response)
	result := rpc.UCI().GetResult(response)
	fmt.Println(reflect.TypeOf(result), result)
}
