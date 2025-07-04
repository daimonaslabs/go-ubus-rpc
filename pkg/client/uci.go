/*
Copyright 2025 Daimonas Labs.

Licensed under the GNU General Public License, Version 3 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.gnu.org/licenses/gpl-3.0.en.html

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"sort"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/dhcp"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/dropbear"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/firewall"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/network"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/system"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/uhttpd"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/wireless"
)

type UCIInterface interface {
	Add(ctx context.Context, opts UCIAddOptions) (r Response, err error)
	Apply(ctx context.Context, opts UCIApplyOptions) (r Response, err error)
	Changes(ctx context.Context, opts UCIChangesOptions) (r Response, err error)
	Configs(ctx context.Context, opts UCIConfigsOptions) (r Response, err error)
	Delete(ctx context.Context, opts UCIDeleteOptions) (r Response, err error)
	Get(ctx context.Context, opts UCIGetOptions) (r Response, err error)
	Revert(ctx context.Context, opts UCIRevertOptions) (r Response, err error)
	Set(ctx context.Context, opts UCISetOptions) (r Response, err error)
}

// implements UCIInterface
type uciRPC struct {
	*UbusRPC
}

func newUCIRPC(u *UbusRPC) *uciRPC {
	u.Call.setPath("uci")
	return &uciRPC{u}
}

func (c *uciRPC) Add(ctx context.Context, opts UCIAddOptions) (Response, error) {
	c.setProcedure("add")
	c.setSignature(opts)

	return c.do(ctx)
}

func (c *uciRPC) Apply(ctx context.Context, opts UCIApplyOptions) (Response, error) {
	c.setProcedure("apply")
	c.setSignature(opts)

	return c.do(ctx)
}

func (c *uciRPC) Changes(ctx context.Context, opts UCIChangesOptions) (Response, error) {
	c.setProcedure("changes")
	c.setSignature(opts)

	return c.do(ctx)
}

func (c *uciRPC) Configs(ctx context.Context, opts UCIConfigsOptions) (Response, error) {
	c.setProcedure("configs")
	c.setSignature(opts)

	return c.do(ctx)
}

func (c *uciRPC) Delete(ctx context.Context, opts UCIDeleteOptions) (Response, error) {
	c.setProcedure("delete")
	c.setSignature(opts)

	return c.do(ctx)
}

func (c *uciRPC) Get(ctx context.Context, opts UCIGetOptions) (Response, error) {
	c.setProcedure("get")
	c.setSignature(opts)

	return c.do(ctx)
}

func (c *uciRPC) Revert(ctx context.Context, opts UCIRevertOptions) (Response, error) {
	c.setProcedure("revert")
	c.setSignature(opts)

	return c.do(ctx)
}

func (c *uciRPC) Set(ctx context.Context, opts UCISetOptions) (Response, error) {
	c.setProcedure("set")
	c.setSignature(opts)

	return c.do(ctx)
}

/*
################################################################
#
# all XOptions types are in this block. they all implement the
# Signature interface.
#
################################################################
*/

// implements Signature interface
type UCIAddOptions struct {
	Config string `json:"config,omitempty"`
	Type   string `json:"type,omitempty"`
}

func (UCIAddOptions) isOptsType() {}

func (opts UCIAddOptions) GetResult(p Response) (u UCIAddResult, err error) {
	if len(p) == 0 {
		return u, errors.New("empty response")
	} else if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case addResult:
			err = json.Unmarshal(data, &u)
		default:
			return u, errors.New("not a UCIAddResult")
		}
	} else { // error
		return u, errors.New(p[0].(ExitCode).Error())
	}
	return u, err
}

// does not have a GetResult func because this command only returns the exit code
// implements Signature interface
type UCIApplyOptions struct {
	Rollback uci.Bool `json:"rollback,omitempty"`
	Timeout  int      `json:"timeout,omitempty"`
}

func (UCIApplyOptions) isOptsType() {}

// implements Signature interface
type UCIChangesOptions struct {
	Config string `json:"config,omitempty"`
}

func (UCIChangesOptions) isOptsType() {}

func (opts UCIChangesOptions) GetResult(p Response) (u UCIChangesResult, err error) {
	u.Changes = make(map[string][]Change)
	if len(p) == 0 {
		return u, errors.New("empty response")
	} else if len(p) > 1 {
		//data, _ := json.Marshal(p[1])
		switch c := p[1].(type) {
		case changesResult:
			if len(c.One) > 0 {
				u.Changes[opts.Config] = exportRawChanges(c.One)
				return u, nil

			} else {
				for config, changes := range c.Many {
					u.Changes[config] = exportRawChanges(changes)
				}
			}
		default:
			return u, errors.New("not a UCIChangesResult")
		}
	} else { // error
		return u, errors.New(p[0].(ExitCode).Error())
	}
	return u, err
}

func exportRawChanges(changes []change) (Changes []Change) {
	for _, c := range changes {
		var C Change
		C.Procedure = c[0]
		C.Section = c[1]
		if len(c) == 3 {
			C.Type = c[2]
		} else if len(c) == 4 {
			C.Option = c[2]
			C.Value = c[3]
		}
		Changes = append(Changes, C)
	}

	return Changes
}

// implements Signature interface
// empty struct because there are no options but it has a special return type so we're
// following the same pattern as the other commands to get the result
type UCIConfigsOptions struct{}

func (UCIConfigsOptions) isOptsType() {}

func (opts UCIConfigsOptions) GetResult(p Response) (u UCIConfigsResult, err error) {

	if len(p) == 0 {
		return u, errors.New("empty response")
	} else if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case configsResult:
			err = json.Unmarshal(data, &u)
		default:
			return u, errors.New("not a UCIConfigsResult")
		}
	} else { // error
		return u, errors.New(p[0].(ExitCode).Error())
	}
	return u, err
}

// does not have a GetResult func because this command only returns the exit code
// implements Signature interface
type UCIDeleteOptions struct {
	Config  string `json:"config,omitempty"`
	Section string `json:"section,omitempty"`
	Type    string `json:"type,omitempty"`
	Option  string `json:"option,omitempty"`
}

func (UCIDeleteOptions) isOptsType() {}

// implements Signature interface
type UCIGetOptions struct {
	Config  string `json:"config,omitempty"`
	Section string `json:"section,omitempty"`
	Type    string `json:"type,omitempty"`
	Option  string `json:"option,omitempty"`
}

func (UCIGetOptions) isOptsType() {}

func (opts UCIGetOptions) GetResult(p Response) (u UCIGetResult, err error) {
	if len(p) == 0 {
		return u, errors.New("empty response")
	} else if len(p) == 1 {
		return u, err
	} else if len(p) > 1 {
		switch obj := p[1].(type) {
		case valueResult:
			u.Option = map[string]uci.List{opts.Option: obj.Value}
		case valuesResult:
			for _, section := range obj.Values {
				switch s := section.(type) {
				case dhcp.BootSection:
					u.Sections = append(u.Sections, s)
				case dhcp.CircuitIDSection:
					u.Sections = append(u.Sections, s)
				case dhcp.DHCPSection:
					u.Sections = append(u.Sections, s)
				case dhcp.DnsmasqSection:
					u.Sections = append(u.Sections, s)
				case dhcp.HostSection:
					u.Sections = append(u.Sections, s)
				case dhcp.HostRecordSection:
					u.Sections = append(u.Sections, s)
				case dhcp.MACSection:
					u.Sections = append(u.Sections, s)
				case dhcp.OdhcpdSection:
					u.Sections = append(u.Sections, s)
				case dhcp.RelaySection:
					u.Sections = append(u.Sections, s)
				case dhcp.RemoteIDSection:
					u.Sections = append(u.Sections, s)
				case dhcp.SubscrIDSection:
					u.Sections = append(u.Sections, s)
				case dhcp.TagSection:
					u.Sections = append(u.Sections, s)
				case dhcp.UserClassSection:
					u.Sections = append(u.Sections, s)
				case dhcp.VendorClassSection:
					u.Sections = append(u.Sections, s)
				case dropbear.DropbearSection:
					u.Sections = append(u.Sections, s)
				case firewall.DefaultsSection:
					u.Sections = append(u.Sections, s)
				case firewall.ForwardingSection:
					u.Sections = append(u.Sections, s)
				case firewall.IPSetSection:
					u.Sections = append(u.Sections, s)
				case firewall.IncludeSection:
					u.Sections = append(u.Sections, s)
				case firewall.RedirectSection:
					u.Sections = append(u.Sections, s)
				case firewall.RuleSection:
					u.Sections = append(u.Sections, s)
				case firewall.ZoneSection:
					u.Sections = append(u.Sections, s)
				case network.BridgeVLANSection:
					u.Sections = append(u.Sections, s)
				case network.DeviceSection:
					u.Sections = append(u.Sections, s)
				case network.GlobalsSection:
					u.Sections = append(u.Sections, s)
				case network.InterfaceSection:
					u.Sections = append(u.Sections, s)
				case network.SwitchSection:
					u.Sections = append(u.Sections, s)
				case network.SwitchPortSection:
					u.Sections = append(u.Sections, s)
				case network.SwitchVLANSection:
					u.Sections = append(u.Sections, s)
				case system.SystemSection:
					u.Sections = append(u.Sections, s)
				case system.TimeserverSection:
					u.Sections = append(u.Sections, s)
				case uhttpd.CertSection:
					u.Sections = append(u.Sections, s)
				case uhttpd.UHTTPdSection:
					u.Sections = append(u.Sections, s)
				case wireless.WifiDeviceSection:
					u.Sections = append(u.Sections, s)
				case wireless.WifiIfaceSection:
					u.Sections = append(u.Sections, s)
				}
			}
		default:
			return u, errors.New("not a UCIGetResult")
		}
	} else { // error
		return u, errors.New(p[0].(ExitCode).Error())
	}
	sort.Slice(u.Sections, func(i, j int) bool {
		return u.Sections[i].GetIndex() < u.Sections[j].GetIndex()
	})
	return u, err
}

// does not have a GetResult func because this command only returns the exit code
// implements Signature interface
type UCIRevertOptions struct {
	Config string `json:"config,omitempty"`
}

func (UCIRevertOptions) isOptsType() {}

// does not have a GetResult func because this command only returns the exit code
// implements Signature interface
type UCISetOptions struct {
	Config  string                   `json:"config,omitempty"`
	Section string                   `json:"section,omitempty"`
	Values  uci.ConfigSectionOptions `json:"values,omitempty"`
}

func (UCISetOptions) isOptsType() {}

/*
################################################################
#
# all exported XResult types are in this block.
#
################################################################
*/

// result of a `uci add` command
type UCIAddResult struct {
	Section string `json:"section,omitempty"`
}

type Change struct {
	Procedure string `json:"procedure"`
	Section   string `json:"section"`
	Type      string `json:"type,omitempty"`
	Option    string `json:"option,omitempty"`
	Value     string `json:"value,omitempty"`
}

type UCIChangesResult struct {
	Changes map[string][]Change `json:"changes"`
}

// result of a `uci configs` command
type UCIConfigsResult struct {
	Configs []string `json:"configs,omitempty"`
}

// result of a `uci get` command
type UCIGetResult struct {
	// if any combination of Config, Section, and Type are specified, return a set of
	// ConfigSection(s)
	Sections []uci.ConfigSection `json:"sections,omitempty"`
	// if Option is set in UCIGetOptions, return a single option's value
	Option map[string]uci.List `json:"option,omitempty"`
}

/*
################################################################
#
# all unexported xResult types are in this block.
#
################################################################
*/

// helper for unmarshaling dynamic xResults objects
type rawMap map[string]json.RawMessage

// implements ResultObject interface
// used for handling the raw RPC response
type addResult struct {
	Section string `json:"section"`
}

func (addResult) isResultObject() {}

type change []string
type changesResult struct {
	Many map[string][]change `json:"many,omitempty"`
	One  []change            `json:"one,omitempty"`
}

func (changesResult) isResultObject() {}

// Many:
//
//	{
//	 "changes": {
//	   "firewall": [
//	     ["set", "cfg04ad58", "enabled", "0"]
//	   ]
//	 }
//	}
//
// One:
//
//	{
//	 "changes": [
//	   ["add", "cfg0fad58", "forwarding"]
//	 ]
//	}
func (v changesResult) MarshalJSON() ([]byte, error) {
	if v.Many != nil {
		manyMap := make(map[string][][]string)
		for section, cmds := range v.Many {
			for _, cmd := range cmds {
				manyMap[section] = append(manyMap[section], []string(cmd))
			}
		}
		return json.Marshal(manyMap)
	}

	if v.One != nil {
		one := make([][]string, len(v.One))
		for i, cmd := range v.One {
			one[i] = []string(cmd)
		}
		return json.Marshal(one)
	}

	return json.Marshal(nil)
}

func (v *changesResult) UnmarshalJSON(data []byte) (err error) {
	// One: [["add", "cfg0fad58", "forwarding" ], ... ] || Many: {"firewall": [["add", "cfg0fad58", "forwarding" ], ... ], "dhcp": [[...], ...]}
	var topLevel rawMap // {"changes": json.RawMessage}
	if err = json.Unmarshal(data, &topLevel); err != nil {
		return err
	}

	changes, ok := topLevel["changes"]

	if !ok {
		return errors.New("malformed changesResult")
	}

	if isSingleChanges(topLevel) {
		err = json.Unmarshal(changes, &v.One)
		if err != nil {
			return err
		}
	} else {
		err = json.Unmarshal(changes, &v.Many)
		if err != nil {
			return err
		}
	}

	return nil
}

func isSingleChanges(m map[string]json.RawMessage) bool {
	var probe any
	if err := json.Unmarshal(m["changes"], &probe); err != nil {
		return false
	}

	return reflect.TypeOf(probe).Kind() == reflect.Slice
}

// implements ResultObject interface
// used for handling the raw RPC response
type configsResult struct {
	Configs []string `json:"configs"`
}

func (configsResult) isResultObject() {}

// implements ResultObject interface
// used for handling the raw RPC response
type valueResult struct {
	Value uci.List `json:"value"`
}

func (valueResult) isResultObject() {}

// implements ResultObject interface
// used for handling the raw RPC response
//
// this struct handles two different types of responses:
// {
//
//	    "values": {
//	        "cfg01e63d": {
//	        	".anonymous": true,
//	         	".type": "defaults",
//	          	".name": "cfg01e63d",
//	          	".index": 0,
//	          	"syn_flood": "1",
//	          	"input": "REJECT",
//	          	"output": "ACCEPT",
//	          	"forward": "REJECT"
//	        },
//	        ...
//		}
//	}
//
// and:
//
// {
//
//	    "values": {
//	        ".anonymous": true,
//	        ".type": "forwarding",
//	        ".name": "cfg04ad58",
//	        "src": "lan",
//	        "dest": "wan"
//	    }
//	}
//
// basically, one is a single object returned while the other is a set of them.
// you can call json.Marshal and Unmarshal on them like normal and it will figure
// out which one it is for you. if it is a single response like in the second example,
// it will unmarshal it into the form of the first one but with only that object.
type valuesResult struct {
	Values map[string]uci.ConfigSection `json:"values"`
}

func (valuesResult) isResultObject() {}

func (v valuesResult) MarshalJSON() ([]byte, error) {
	// handle single unnamed config (e.g., ".type": "forwarding")
	if len(v.Values) == 1 {
		for _, section := range v.Values {
			return json.Marshal(section)
		}
	}

	// otherwise marshal as map
	out := make(map[string]interface{})
	for name, section := range v.Values {
		out[name] = section
	}
	return json.Marshal(out)
}

func (v *valuesResult) UnmarshalJSON(data []byte) (err error) {
	var topLevel rawMap // {"values": json.RawMessage}
	var result rawMap   // SINGLE: {".anonymous": json.RawMessage} || MULTIPLE: {"cfg04ad58": json.RawMessage}
	if err := json.Unmarshal(data, &topLevel); err != nil {
		return err
	}

	values, ok := topLevel["values"]

	if !ok {
		return errors.New("malformed valuesResult")
	}

	if err := json.Unmarshal(values, &result); err != nil {
		return err
	}

	if isSingleValues(result) {
		section, err := unmarshalRawSection(values)
		if err != nil {
			return err
		}
		v.Values = map[string]uci.ConfigSection{section.GetName(): section}
		return nil
	} else {
		// handle named entries in map
		v.Values = make(map[string]uci.ConfigSection)
		for name, section := range result {
			section, err := unmarshalRawSection(section)
			if err != nil {
				return err
			}
			v.Values[name] = section
		}
	}

	return nil
}

func unmarshalRawResult[S uci.ConfigSection](data []byte) (uci.ConfigSection, error) {
	var s S
	err := json.Unmarshal(data, &s)
	return s, err
}

func unmarshalRawSection(data []byte) (section uci.ConfigSection, err error) {
	var probe struct {
		Type string `json:".type"`
	}
	var rawSection rawMap

	if err = json.Unmarshal(data, &rawSection); err != nil {
		return nil, err
	}

	sectionBytes, err := json.Marshal(rawSection)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(sectionBytes, &probe); err != nil {
		return nil, err
	}

	switch probe.Type {
	case string(dhcp.Boot):
		section, err = unmarshalRawResult[dhcp.BootSection](data)
	case string(dhcp.CircuitID):
		section, err = unmarshalRawResult[dhcp.CircuitIDSection](data)
	case string(dhcp.DHCP):
		section, err = unmarshalRawResult[dhcp.DHCPSection](data)
	case string(dhcp.Dnsmasq):
		section, err = unmarshalRawResult[dhcp.DnsmasqSection](data)
	case string(dhcp.Host):
		section, err = unmarshalRawResult[dhcp.HostSection](data)
	case string(dhcp.HostRecord):
		section, err = unmarshalRawResult[dhcp.HostRecordSection](data)
	case string(dhcp.MAC):
		section, err = unmarshalRawResult[dhcp.MACSection](data)
	case string(dhcp.Odhcpd):
		section, err = unmarshalRawResult[dhcp.OdhcpdSection](data)
	case string(dhcp.Relay):
		section, err = unmarshalRawResult[dhcp.RelaySection](data)
	case string(dhcp.RemoteID):
		section, err = unmarshalRawResult[dhcp.RemoteIDSection](data)
	case string(dhcp.SubscrID):
		section, err = unmarshalRawResult[dhcp.SubscrIDSection](data)
	case string(dhcp.Tag):
		section, err = unmarshalRawResult[dhcp.TagSection](data)
	case string(dhcp.UserClass):
		section, err = unmarshalRawResult[dhcp.UserClassSection](data)
	case string(dhcp.VendorClass):
		section, err = unmarshalRawResult[dhcp.VendorClassSection](data)
	case string(dropbear.Dropbear):
		section, err = unmarshalRawResult[dropbear.DropbearSection](data)
	case string(firewall.Defaults):
		section, err = unmarshalRawResult[firewall.DefaultsSection](data)
	case string(firewall.Forwarding):
		section, err = unmarshalRawResult[firewall.ForwardingSection](data)
	case string(firewall.IPSet):
		section, err = unmarshalRawResult[firewall.IPSetSection](data)
	case string(firewall.Include):
		section, err = unmarshalRawResult[firewall.IncludeSection](data)
	case string(firewall.Redirect):
		section, err = unmarshalRawResult[firewall.RedirectSection](data)
	case string(firewall.Rule):
		section, err = unmarshalRawResult[firewall.RuleSection](data)
	case string(firewall.Zone):
		section, err = unmarshalRawResult[firewall.ZoneSection](data)
	case string(network.BridgeVLAN):
		section, err = unmarshalRawResult[network.BridgeVLANSection](data)
	case string(network.Device):
		section, err = unmarshalRawResult[network.DeviceSection](data)
	case string(network.Globals):
		section, err = unmarshalRawResult[network.GlobalsSection](data)
	case string(network.Interface):
		section, err = unmarshalRawResult[network.InterfaceSection](data)
	case string(network.Switch):
		section, err = unmarshalRawResult[network.SwitchSection](data)
	case string(network.SwitchPort):
		section, err = unmarshalRawResult[network.SwitchPortSection](data)
	case string(network.SwitchVLAN):
		section, err = unmarshalRawResult[network.SwitchVLANSection](data)
	case string(system.System):
		section, err = unmarshalRawResult[system.SystemSection](data)
	case string(system.Timeserver):
		section, err = unmarshalRawResult[system.TimeserverSection](data)
	case string(uhttpd.Cert):
		section, err = unmarshalRawResult[uhttpd.CertSection](data)
	case string(uhttpd.UHTTPd):
		section, err = unmarshalRawResult[uhttpd.UHTTPdSection](data)
	case string(wireless.WifiDevice):
		section, err = unmarshalRawResult[wireless.WifiDeviceSection](data)
	case string(wireless.WifiIface):
		section, err = unmarshalRawResult[wireless.WifiIfaceSection](data)
	default:
		return nil, errors.New("invalid config section")
	}
	return section, err
}

// checks if the value of `values` is a single uci.ConfigSection or not
func isSingleValues(m map[string]json.RawMessage) bool {
	_, ok := m[".anonymous"]
	return ok
}

/*
################################################################
#
# all matchXResult funcs are in this block. used in init().
#
################################################################
*/

// matcher for addResult
func matchAddResult(data json.RawMessage) (ResultObject, error) {
	var val addResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Section) > 0 {
			return val, nil
		}
	}

	return nil, nil
}

// matcher for changesResult
func matchChangesResult(data json.RawMessage) (ResultObject, error) {
	var raw rawMap
	var val changesResult

	if err := json.Unmarshal(data, &raw); err == nil {
		if _, ok := raw["changes"]; ok {
			err = json.Unmarshal(data, &val)
			if err == nil {
				return val, err
			}
		}
	}

	return nil, nil
}

// matcher for configsResult
func matchConfigsResult(data json.RawMessage) (ResultObject, error) {
	var val configsResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Configs) > 0 {
			return val, nil
		}
	}

	return nil, nil
}

// matcher for valueResult
func matchValueResult(data json.RawMessage) (ResultObject, error) {
	var val valueResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Value) > 0 {
			return val, nil
		}
	}

	return nil, nil
}

// matcher for valuesResult
func matchValuesResult(data json.RawMessage) (ResultObject, error) {
	var val valuesResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Values) > 0 {
			return val, nil
		}
	}

	return nil, nil
}
