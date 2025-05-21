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
	uciGetOpts := client.UCIGetOptions{Config: "firewall", Section: "cfg0b92bd"} //, Option: "icmp_type"}
	//sessionLoginOpts := client.SessionLoginOptions{Username: "root", Password: "D@!monas"}
	//rpc.Call = rpc.Session().Login(&sessionLoginOpts)
	//rpc.Call = rpc.UCI().Configs()
	response, err := rpc.UCI().Get(ctx, uciGetOpts)
	if err != nil {
		fmt.Println("main1: ", err)
	}
	result, err := uciGetOpts.GetResult(response)
	//result, err := sessionLoginOpts.GetResult(response)
	//result, err := client.UCIConfigsOptions{}.GetResult(response)
	fmt.Println("response: ", response)
	fmt.Println("result: ", result)
	//fmt.Println("result: ", reflect.TypeOf(result), result.SectionArray)
	//	for i, s := range result.SectionArray {
	//fmt.Println("Go index: ", i, " ubus index: ", s.GetIndex())
	//}
	fmt.Println("err: ", reflect.TypeOf(err), err)
}

/*
curl -k -H 'Content-Type: application/json' -d '{ "jsonrpc": "2.0", "id": 1, "method": "call", "params": [ "9a755e11de0c51431c6f1bc5ed0b885e", "uci", "set", {"config": "firewall", "section":"cfg04ad58", "values":{"enabled":true} } ] }'  https://10.0.0.1/ubus | jq -r
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": [
    0
  ]
}

`uci set` only works on existing sections and errors out if you provide the static fields in the `values` portion

*/
