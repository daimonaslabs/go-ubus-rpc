package client

import (
	"context"
	"flag"
	"log"
	"reflect"
	"slices"
	"testing"

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
	opts := ClientOptions{Username: *username, Password: *password, URL: *url, Timeout: 15}
	rpc, err := NewUbusRPC(ctx, &opts)
	if err != nil {
		log.Fatalln("error creating ubus client")
	}

	return ctx, rpc
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func TestUCIAddSetDelete(t *testing.T) {
	ctx, rpc := prepare()

	// add a new config section and set an option within it
	uciAddOpts := UCIAddOptions{Config: firewall.Config, Type: firewall.Forwarding}
	addResponse, err := rpc.UCI().Add(ctx, uciAddOpts)
	checkErr(t, err)
	addResult, err := uciAddOpts.GetResult(addResponse)
	checkErr(t, err)

	forwardingSectionOptions := firewall.ForwardingSectionOptions{
		Enabled: uci.StringBoolTrue,
	}
	uciSetOpts := UCISetOptions{Config: firewall.Config, Section: addResult.Section, Values: forwardingSectionOptions}
	_, err = rpc.UCI().Set(ctx, uciSetOpts)
	checkErr(t, err)

	uciApplyOpts := UCIApplyOptions{Rollback: uci.StringBoolTrue, Timeout: 10}
	_, err = rpc.UCI().Apply(ctx, uciApplyOpts)
	checkErr(t, err)

	// check that the config was actually applied
	uciGetOpts := UCIGetOptions{Config: firewall.Config, Section: addResult.Section}
	getResponse, err := rpc.UCI().Get(ctx, uciGetOpts)
	checkErr(t, err)
	getResult, err := uciGetOpts.GetResult(getResponse)
	checkErr(t, err)

	newSection, ok := getResult.SectionArray[0].(firewall.ForwardingSection)
	if !ok {
		t.Error("result is not a ForwardingSection")
	}

	if !reflect.DeepEqual(newSection.ForwardingSectionOptions, forwardingSectionOptions) {
		t.Error("options not set correctly")
	}
	t.Log("\nexpected result: ", forwardingSectionOptions, "\nactual result: ", newSection.ForwardingSectionOptions)

	// delete the section
	uciDeleteOpts := UCIDeleteOptions{Config: firewall.Config, Section: addResult.Section}
	_, err = rpc.UCI().Delete(ctx, uciDeleteOpts)
	checkErr(t, err)
	_, err = rpc.UCI().Apply(ctx, uciApplyOpts)
	checkErr(t, err)

	// confirm deletion
	_, err = rpc.UCI().Get(ctx, uciGetOpts)
	if err != nil {
		t.Error(err)
	}
}

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
		t.Error("\nexpected configsResult: ", expected.Configs, "\nactual configsResult: ", result.Configs)
	} else {
		t.Log("\nexpected configsResult: ", expected.Configs, "\nactual configsResult: ", result.Configs, "\nactual configsResult contains all expected values")
	}
}

func TestUCIRevert(t *testing.T) {
	ctx, rpc := prepare()
	uciAddOpts := UCIAddOptions{Config: firewall.Config, Type: firewall.Forwarding}
	uciChangesOpts := UCIChangesOptions{Config: firewall.Config}
	uciRevertOpts := UCIRevertOptions{Config: firewall.Config}

	rpc.UCI().Add(ctx, uciAddOpts)
	changesResponse, _ := rpc.UCI().Changes(ctx, uciChangesOpts)
	changesResult, _ := uciChangesOpts.GetResult(changesResponse)

	t.Log("should be one change: ", changesResult)
	if changesResult.Changes == nil {
		t.Error("problem listing changes")
	}

	revertResponse, err := rpc.UCI().Revert(ctx, uciRevertOpts)
	t.Log("should be exit status 0: ", revertResponse)
	if err != nil {
		t.Error(err)
	}
	changesResponse, _ = rpc.UCI().Changes(ctx, uciChangesOpts)
	changesResult, _ = uciChangesOpts.GetResult(changesResponse)
	t.Log("should be zero changes: ", changesResult)
	if changesResult.Changes == nil {
		t.Error("did not revert changes!")
	}
}
