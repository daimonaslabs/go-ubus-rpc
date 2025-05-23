package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/firewall"
)

func main() {
	ctx := context.TODO()
	opts := client.ClientOptions{Username: "root", Password: "D@!monas", URL: "http://10.0.0.1/ubus", Timeout: session.DefaultSessionTimeout}
	rpc, err := client.NewUbusRPC(ctx, &opts)

	if err != nil {
		fmt.Println(rpc, err)
	}
	uciAddOpts := client.UCIAddOptions{Config: firewall.Config, Type: firewall.Forwarding}
	response, _ := rpc.UCI().Add(ctx, uciAddOpts)
	result, err := uciAddOpts.GetResult(response)

	//uciGetOpts := client.UCIGetOptions{Config: "firewall", Section: "cfg0b92bd"} //, Option: "icmp_type"}
	//response, err := rpc.UCI().Get(ctx, uciGetOpts)
	//result, err := uciGetOpts.GetResult(response)

	//sessionLoginOpts := client.SessionLoginOptions{Username: "root", Password: "D@!monas"}
	//response, err := rpc.UCI().Session.Login(ctx, uciGetOpts)
	//result, err := sessionLoginOpts.GetResult(response)

	////forwarding := firewall.ForwardingSectionOptions{
	////	Enabled: uci.StringBoolTrue,
	////}
	//uciSetOpts := client.UCISetOptions{Config: firewall.Config, Section: "cfg04ad58", Values: forwarding}
	//response, _ := rpc.UCI().Set(ctx, uciSetOpts)
	//fmt.Println("set response: ", response)

	uciApplyOpts := client.UCIApplyOptions{Rollback: uci.StringBoolTrue, Timeout: 10}
	response, err = rpc.UCI().Apply(ctx, uciApplyOpts)

	//fmt.Println("result: ", reflect.TypeOf(result), result.SectionArray)
	//	for i, s := range result.SectionArray {
	//fmt.Println("Go index: ", i, " ubus index: ", s.GetIndex())
	//}
	fmt.Println("result: ", result)
	fmt.Println("err: ", reflect.TypeOf(err), err)
}
