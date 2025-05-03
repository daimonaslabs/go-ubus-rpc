package main

import (
	"context"
	"fmt"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
)

func main() {
	ctx := context.TODO()
	opts := client.ClientOptions{Username: "root", Password: "D@!monas", URL: "http://10.0.0.1/ubus", Timeout: uint(5)}
	rpc, err := client.NewUbusRPC(ctx, &opts)

	if err != nil {
		fmt.Println(rpc, err)
	}

	rpc.Call = rpc.UCI().Get(&client.UCIOptions{Config: "firewall", Section: "cfg04ad58", Option: "src"})
	result, err := rpc.Do(ctx)
	data := new([]byte)
	result.UnmarshalJSON(*data)
	fmt.Println(data)
}

// SAMPLE CALL BY USER
//	login := client.LoginOptions{
//		Username: opts.Username,
//		Password: opts.Password,
//		Timeout:  client.DefaultSessionTimeout,
//	}
//	rpc.Call = rpc.Session().Login(&login)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	response, err := rpc.Do(ctx)
