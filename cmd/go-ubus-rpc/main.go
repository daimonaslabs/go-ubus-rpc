package main

import (
	"context"
	"fmt"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
)

func main() {
	ctx := context.TODO()
	opts := client.ClientOptions{Username: "root", Password: "D@!monas", URL: "http://10.0.0.1/ubus", Timeout: client.DefaultSessionTimeout}
	rpc, err := client.NewUbusRPC(ctx, &opts)

	if err != nil {
		fmt.Println(rpc, err)
	}

	rpc.Call = rpc.UCI().Get(&client.UCIOptions{Config: "firewall", Section: "cfg04ad58", Option: "src"})
	fmt.Println("main Call before Do(): ", rpc.Call)
	response, err := rpc.Do(ctx)
	result, _ := client.GetAs[client.ValueResult](response)
	fmt.Println("main Result after Do(): ", result)
}
