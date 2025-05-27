package main

import (
	"context"
	"fmt"

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
	uciAddOpts := client.UCIAddOptions{Config: firewall.Config, Type: firewall.Forwarding}
	response, _ := rpc.UCI().Add(ctx, uciAddOpts)
	//addResult, err := uciAddOpts.GetResult(response)
	fmt.Println("addResponse: ", response, err)
	uciChangesOpts := client.UCIChangesOptions{} //Config: firewall.Config}
	response, _ = rpc.UCI().Changes(ctx, uciChangesOpts)
	fmt.Println("addChangesResponse: ", response, err)
	//uciApplyOpts := client.UCIApplyOptions{Rollback: uci.StringBoolTrue, Timeout: 10}
	//response, err = rpc.UCI().Apply(ctx, uciApplyOpts)
	//fmt.Println("addApplyResponse: ", response, err)

	uciRevertOpts := client.UCIRevertOptions{Config: firewall.Config}
	response, _ = rpc.UCI().Revert(ctx, uciRevertOpts)
	fmt.Println("revertResponse: ", response, err)
	response, _ = rpc.UCI().Changes(ctx, uciChangesOpts)
	fmt.Println("addChangesResponse: ", response, err)

	//forwarding := firewall.ForwardingSectionOptions{
	//	Enabled: uci.StringBoolTrue,
	//}
	//uciSetOpts := client.UCISetOptions{Config: firewall.Config, Section: addResult.Section, Values: forwarding}
	//response, err = rpc.UCI().Set(ctx, uciSetOpts)
	//fmt.Println("setResponse: ", response, err)
	//response, _ = rpc.UCI().Changes(ctx, uciChangesOpts)
	//fmt.Println("setChangesResponse: ", response, err)
	//response, err = rpc.UCI().Apply(ctx, uciApplyOpts)
	//fmt.Println("setApplyResponse: ", response, err)

	//uciDelOpts := client.UCIDeleteOptions{Config: firewall.Config, Section: addResult.Section} //, Option: "enabled"}
	//response, err = rpc.UCI().Delete(ctx, uciDelOpts)
	//fmt.Println("delResponse: ", response, err)
	//response, _ = rpc.UCI().Changes(ctx, uciChangesOpts)
	//fmt.Println("delChangesResponse: ", response, err)
	//response, err = rpc.UCI().Apply(ctx, uciApplyOpts)
	//fmt.Println("delApplyResponse: ", response, err)

	//uciAddOpts2 := client.UCIAddOptions{Config: firewall.Config, Type: firewall.Zone}
	//_, _ = rpc.UCI().Add(ctx, uciAddOpts2)
	//result, err := uciAddOpts.GetResult(response)
	//fmt.Println("addResponse: ", response)
	//fmt.Println("addResult: ", result)
	//fmt.Println("addErr: ", err)

	//uciChangesOpts := client.UCIChangesOptions{} //Config: firewall.Config}
	//response, _ = rpc.UCI().Changes(ctx, uciChangesOpts)
	//changesResult, err := uciChangesOpts.GetResult(response)

	//fmt.Println("changesResponse: ", response)
	//fmt.Println("changesResult: ", changesResult)
	//fmt.Println("changesErr: ", err)

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

	//uciApplyOpts := client.UCIApplyOptions{Rollback: uci.StringBoolTrue, Timeout: 10}
	//response, err = rpc.UCI().Apply(ctx, uciApplyOpts)

	//fmt.Println("result: ", reflect.TypeOf(result), result.SectionArray)
	//	for i, s := range result.SectionArray {
	//fmt.Println("Go index: ", i, " ubus index: ", s.GetIndex())
	//}
}
