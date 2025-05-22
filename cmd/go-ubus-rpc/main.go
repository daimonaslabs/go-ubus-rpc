package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/firewall"
)

func main() {
	ctx := context.TODO()
	opts := client.ClientOptions{Username: "root", Password: "D@!monas", URL: "http://10.0.0.1/ubus", Timeout: session.DefaultSessionTimeout}
	rpc, err := client.NewUbusRPC(ctx, &opts)

	if err != nil {
		fmt.Println(rpc, err)
	}
	//uciGetOpts := client.UCIGetOptions{Config: "firewall", Section: "cfg0b92bd"} //, Option: "icmp_type"}
	//response, err := rpc.UCI().Get(ctx, uciGetOpts)
	//result, err := uciGetOpts.GetResult(response)

	//sessionLoginOpts := client.SessionLoginOptions{Username: "root", Password: "D@!monas"}
	//response, err := rpc.UCI().Session.Login(ctx, uciGetOpts)
	//result, err := sessionLoginOpts.GetResult(response)
	forwarding := firewall.ForwardingSectionOptions{
		Enabled: "hello?",
		Family:  "ipv4",
	}
	uciSetOpts := client.UCISetOptions{Config: firewall.Config, Section: "cfg04ad58", Values: forwarding}
	if err != nil {
		fmt.Println("main1: ", err)
	}
	response, err := rpc.UCI().Set(ctx, uciSetOpts)
	fmt.Println("response: ", response)
	//fmt.Println("result: ", result)
	//fmt.Println("result: ", reflect.TypeOf(result), result.SectionArray)
	//	for i, s := range result.SectionArray {
	//fmt.Println("Go index: ", i, " ubus index: ", s.GetIndex())
	//}
	fmt.Println("err: ", reflect.TypeOf(err), err)
}
