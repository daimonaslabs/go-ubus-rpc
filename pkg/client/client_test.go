package client

import (
	"context"
	"flag"
	"log"
	"slices"
	"testing"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/firewall"
)

var (
	username = flag.String("username", "root", "Username to log into OpenWRT instance")
	password = flag.String("password", "D@!monas", "Password to log into OpenWRT instance")
	url      = flag.String("url", "http://10.0.0.1/ubus", "URL of ubus endpoint")
)

func prepare() (ctx context.Context, rpc *UbusRPC) {
	ctx = context.Background()
	opts := ClientOptions{Username: *username, Password: *password, URL: *url, Timeout: session.DefaultSessionTimeout}
	rpc, err := NewUbusRPC(ctx, &opts)
	if err != nil {
		log.Fatalln("error creating ubus client")
	}

	return ctx, rpc
}

/*
	func TestMain(m *testing.M) {
		flag.Parse()

			sequence1:
			get entire config
			add
			set
			$reusable

			sequence2:
			get entire config
			delete
			$reusable

			reusable:
			changes
			apply
			get entire config

			independent:
			configs
			revert

		ctx, rpc := prepare()

		uciAddOpts := UCIAddOptions{Config: firewall.Config, Type: firewall.Forwarding}
		response, _ := rpc.UCI().Add(ctx, uciAddOpts)
		addResult, _ := uciAddOpts.GetResult(response)

		uciChangesOpts := UCIChangesOptions{} //Config: firewall.Config}
		response, _ = rpc.UCI().Changes(ctx, uciChangesOpts)

		uciApplyOpts := UCIApplyOptions{Rollback: uci.StringBoolTrue, Timeout: 10}
		response, _ = rpc.UCI().Apply(ctx, uciApplyOpts)

		uciRevertOpts := UCIRevertOptions{Config: firewall.Config}
		response, _ = rpc.UCI().Revert(ctx, uciRevertOpts)

		response, _ = rpc.UCI().Changes(ctx, uciChangesOpts)

		forwarding := firewall.ForwardingSectionOptions{
			Enabled: uci.StringBoolTrue,
		}
		uciSetOpts := UCISetOptions{Config: firewall.Config, Section: addResult.Section, Values: forwarding}
		response, _ = rpc.UCI().Set(ctx, uciSetOpts)

		response, _ = rpc.UCI().Changes(ctx, uciChangesOpts)
		response, _ = rpc.UCI().Apply(ctx, uciApplyOpts)

		uciDelOpts := UCIDeleteOptions{Config: firewall.Config, Section: addResult.Section} //, Option: "enabled"}
		response, _ = rpc.UCI().Delete(ctx, uciDelOpts)

		response, _ = rpc.UCI().Changes(ctx, uciChangesOpts)

		response, _ = rpc.UCI().Apply(ctx, uciApplyOpts)

		uciAddOpts2 := UCIAddOptions{Config: firewall.Config, Type: firewall.Zone}
		_, _ = rpc.UCI().Add(ctx, uciAddOpts2)
		//result, err := uciAddOpts.GetResult(response)

		response, _ = rpc.UCI().Changes(ctx, uciChangesOpts)
		//changesResult, err := uciChangesOpts.GetResult(response)

		uciGetOpts := UCIGetOptions{Config: "firewall", Section: "cfg0b92bd"} //, Option: "icmp_type"}
		response, _ = rpc.UCI().Get(ctx, uciGetOpts)
		//getResult, err := uciGetOpts.GetResult(response)

		sessionLoginOpts := SessionLoginOptions{Username: "root", Password: "D@!monas"}
		response, _ = rpc.Session().Login(ctx, sessionLoginOpts)
		//loginResult, err := sessionLoginOpts.GetResult(response)

		forwarding = firewall.ForwardingSectionOptions{
			Enabled: uci.StringBoolTrue,
		}
		uciSetOpts = UCISetOptions{Config: firewall.Config, Section: "cfg04ad58", Values: forwarding}
		response, _ = rpc.UCI().Set(ctx, uciSetOpts)

		uciApplyOpts = UCIApplyOptions{Rollback: uci.StringBoolTrue, Timeout: 10}
		response, _ = rpc.UCI().Apply(ctx, uciApplyOpts)

}
*/
func TestUCIConfigs(t *testing.T) {
	ctx, rpc := prepare()
	expected := configsResult{Configs: []uci.ConfigName{uci.DHCP, uci.Dropbear, uci.Firewall, uci.LuCI,
		uci.Network, uci.RPCD, uci.System, uci.UBootEnv, uci.UCITrack, uci.UHTTPd, uci.Wireless}}

	uciConfigsOpts := UCIConfigsOptions{}
	response, err := rpc.UCI().Configs(ctx, uciConfigsOpts)
	if err != nil {
		t.Error(err)
	}
	result, err := uciConfigsOpts.GetResult(response)
	if err != nil {
		t.Error(err)
	}

	var notPresent bool
	for _, config := range expected.Configs {
		if !slices.Contains(result.Configs, config) {
			notPresent = true
		}
	}
	if notPresent {
		t.Error("expected configsResult: ", expected.Configs)
		t.Error("acctual configsResult: ", result.Configs)
	}
}

func TestUCIRevert(t *testing.T) {
	ctx, rpc := prepare()
	expected := changesResult{}
	uciRevertOpts := UCIRevertOptions{Config: firewall.Config}

}
