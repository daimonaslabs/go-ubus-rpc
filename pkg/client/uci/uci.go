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

package uci

import (
	"context"
	"encoding/json"
	"errors"
	"sort"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/rpc"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus"
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
	Add(ctx context.Context, opts AddOptions) (r ubus.Response, err error)
	Apply(ctx context.Context, opts ApplyOptions) (r ubus.Response, err error)
	Changes(ctx context.Context, opts ChangesOptions) (r ubus.Response, err error)
	Configs(ctx context.Context, opts ConfigsOptions) (r ubus.Response, err error)
	Delete(ctx context.Context, opts DeleteOptions) (r ubus.Response, err error)
	Get(ctx context.Context, opts GetOptions) (r ubus.Response, err error)
	Revert(ctx context.Context, opts RevertOptions) (r ubus.Response, err error)
	Set(ctx context.Context, opts SetOptions) (r ubus.Response, err error)
}

// implements UCIInterface
type UCIClient struct {
	call   *ubus.Call
	client rpc.UbusRPCInterface
}

func NewUCIClient(r *rpc.UbusRPCClient) (c *UCIClient) {
	c = &UCIClient{
		call: &ubus.Call{},
	}
	c.call.SetPath("uci")
	c.client = r
	return c
}

func (c *UCIClient) Add(ctx context.Context, opts AddOptions) (ubus.Response, error) {
	c.call.SetProcedure("add")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

func (c *UCIClient) Apply(ctx context.Context, opts ApplyOptions) (ubus.Response, error) {
	c.call.SetProcedure("apply")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

func (c *UCIClient) Changes(ctx context.Context, opts ChangesOptions) (ubus.Response, error) {
	c.call.SetProcedure("changes")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

func (c *UCIClient) Configs(ctx context.Context, opts ConfigsOptions) (ubus.Response, error) {
	c.call.SetProcedure("configs")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

func (c *UCIClient) Delete(ctx context.Context, opts DeleteOptions) (ubus.Response, error) {
	c.call.SetProcedure("delete")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

func (c *UCIClient) Get(ctx context.Context, opts GetOptions) (ubus.Response, error) {
	c.call.SetProcedure("get")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

func (c *UCIClient) Revert(ctx context.Context, opts RevertOptions) (ubus.Response, error) {
	c.call.SetProcedure("revert")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

func (c *UCIClient) Set(ctx context.Context, opts SetOptions) (ubus.Response, error) {
	c.call.SetProcedure("set")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
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
type AddOptions struct {
	Config string `json:"config,omitempty"`
	Type   string `json:"type,omitempty"`
}

func (AddOptions) IsOptsType() {}

func (opts AddOptions) GetResult(p ubus.Response) (u AddResult, err error) {
	if len(p) == 0 {
		return u, errors.New("empty response")
	} else if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case ubus.AddResult:
			err = json.Unmarshal(data, &u)
		default:
			return u, errors.New("not a UCIAddResult")
		}
	} else { // error
		return u, errors.New(p[0].(ubus.ExitCode).Error())
	}
	return u, err
}

// does not have a GetResult func because this command only returns the exit code
// implements Signature interface
type ApplyOptions struct {
	Rollback uci.Bool `json:"rollback,omitempty"`
	Timeout  int      `json:"timeout,omitempty"`
}

func (ApplyOptions) IsOptsType() {}

// implements Signature interface
type ChangesOptions struct {
	Config string `json:"config,omitempty"`
}

func (ChangesOptions) IsOptsType() {}

func (opts ChangesOptions) GetResult(p ubus.Response) (u ChangesResult, err error) {
	u.Changes = make(map[string][]Change)
	if len(p) == 0 {
		return u, errors.New("empty response")
	} else if len(p) > 1 {
		//data, _ := json.Marshal(p[1])
		switch c := p[1].(type) {
		case ubus.ChangesResult:
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
		return u, errors.New(p[0].(ubus.ExitCode).Error())
	}
	return u, err
}

func exportRawChanges(changes []ubus.Change) (Changes []Change) {
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
type ConfigsOptions struct{}

func (ConfigsOptions) IsOptsType() {}

func (opts ConfigsOptions) GetResult(p ubus.Response) (u ConfigsResult, err error) {

	if len(p) == 0 {
		return u, errors.New("empty response")
	} else if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case ubus.ConfigsResult:
			err = json.Unmarshal(data, &u)
		default:
			return u, errors.New("not a UCIConfigsResult")
		}
	} else { // error
		return u, errors.New(p[0].(ubus.ExitCode).Error())
	}
	return u, err
}

// does not have a GetResult func because this command only returns the exit code
// implements Signature interface
type DeleteOptions struct {
	Config  string `json:"config,omitempty"`
	Section string `json:"section,omitempty"`
	Type    string `json:"type,omitempty"`
	Option  string `json:"option,omitempty"`
}

func (DeleteOptions) IsOptsType() {}

// implements Signature interface
type GetOptions struct {
	Config  string `json:"config,omitempty"`
	Section string `json:"section,omitempty"`
	Type    string `json:"type,omitempty"`
	Option  string `json:"option,omitempty"`
}

func (GetOptions) IsOptsType() {}

func (opts GetOptions) GetResult(p ubus.Response) (u GetResult, err error) {
	if len(p) == 0 {
		return u, errors.New("empty response")
	} else if len(p) == 1 {
		return u, err
	} else if len(p) > 1 {
		switch obj := p[1].(type) {
		case ubus.ValueResult:
			u.Option = map[string]uci.List{opts.Option: obj.Value}
		case ubus.ValuesResult:
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
		return u, errors.New(p[0].(ubus.ExitCode).Error())
	}
	sort.Slice(u.Sections, func(i, j int) bool {
		return u.Sections[i].GetIndex() < u.Sections[j].GetIndex()
	})
	return u, err
}

// does not have a GetResult func because this command only returns the exit code
// implements Signature interface
type RevertOptions struct {
	Config string `json:"config,omitempty"`
}

func (RevertOptions) IsOptsType() {}

// does not have a GetResult func because this command only returns the exit code
// implements Signature interface
type SetOptions struct {
	Config  string                   `json:"config,omitempty"`
	Section string                   `json:"section,omitempty"`
	Values  uci.ConfigSectionOptions `json:"values,omitempty"`
}

func (SetOptions) IsOptsType() {}

/*
################################################################
#
# all exported XResult types are in this block.
#
################################################################
*/

// result of a `uci add` command
type AddResult struct {
	Section string `json:"section,omitempty"`
}

type Change struct {
	Procedure string `json:"procedure"`
	Section   string `json:"section"`
	Type      string `json:"type,omitempty"`
	Option    string `json:"option,omitempty"`
	Value     string `json:"value,omitempty"`
}

type ChangesResult struct {
	Changes map[string][]Change `json:"changes"`
}

// result of a `uci configs` command
type ConfigsResult struct {
	Configs []string `json:"configs,omitempty"`
}

// result of a `uci get` command
type GetResult struct {
	// if any combination of Config, Section, and Type are specified, return a set of
	// ConfigSection(s)
	Sections []uci.ConfigSection `json:"sections,omitempty"`
	// if Option is set in UCIGetOptions, return a single option's value
	Option map[string]uci.List `json:"option,omitempty"`
}
