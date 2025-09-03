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

package ubus

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/dhcp"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/dropbear"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/firewall"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/network"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/system"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/uhttpd"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/wireless"
)

// interface for content within a Response
// all ResultObjects must also have their own match function
type ResultObject interface {
	IsResultObject()
}

// implements json.Marshaler and json.Unmarshaler
// effectively a tuple:
// Response[0] is always an int (ExitCode)
// Response[1] is always an xResult type (e.g. UCIResult)
type Response []ResultObject

// custom UnmarshalJSON for Response
func (r *Response) UnmarshalJSON(data []byte) error {
	var rawLen, matches int
	var raw []json.RawMessage

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	rawLen = len(raw)

	for _, item := range raw {
		var matched bool
		for _, matcher := range resultObjectMatcherRegistry {
			if obj, err := matcher(item); err == nil && obj != nil {
				*r = append(*r, obj)
				matches += 1
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf("unknown result object: %s", string(item))
		}
	}

	if matches != rawLen {
		return fmt.Errorf("error parsing Response object")
	}

	return nil
}

// custom MarshalJSON for Response
func (r Response) MarshalJSON() ([]byte, error) {
	var raw []json.RawMessage

	for _, obj := range r {
		data, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}
		raw = append(raw, data)
	}

	return json.Marshal(raw)
}

/*
################################################################
#
# all raw ubus respose XResult types are in this block.
#
################################################################
*/

// implements ResultObject
// implements builtin.error
// always the first object of the Response tuple
type ExitCode int

func (e ExitCode) IsResultObject() {}

const (
	UbusStatusOK ExitCode = iota
	UbusStatusInvalidCommand
	UbusStatusInvalidArgument
	UbusStatusMethodNotFound
	UbusStatusNotFound
	UbusStatusNoData
	UbusStatusPermissionDenied
	UbusStatusTimeout
	UbusStatusNotSupported
	UbusStatusUnknownError
	UbusStatusConnectionFailed
	UbusStatusLast
)

func (e ExitCode) String() string {
	switch e {
	case UbusStatusOK:
		return "UbusStatusOK"
	case UbusStatusInvalidCommand:
		return "UbusStatusInvalidCommand"
	case UbusStatusInvalidArgument:
		return "UbusStatusInvalidArgument"
	case UbusStatusMethodNotFound:
		return "UbusStatusMethodNotFound"
	case UbusStatusNotFound:
		return "UbusStatusNotFound"
	case UbusStatusNoData:
		return "UbusStatusNoData"
	case UbusStatusPermissionDenied:
		return "UbusStatusPermissionDenied"
	case UbusStatusTimeout:
		return "UbusStatusTimeout"
	case UbusStatusNotSupported:
		return "UbusStatusNotSupported"
	case UbusStatusUnknownError:
		return "UbusStatusUnknownError"
	case UbusStatusConnectionFailed:
		return "UbusStatusConnectionFailed"
	case UbusStatusLast:
		return "UbusStatusLast"
	default:
		return "Unknown ExitCode"
	}
}

func (e ExitCode) Error() string {
	return fmt.Sprintf("exit status (%d) %s", e, e.String())
}

// implements ResultObject interface
// used for handling the raw RPC response
type SessionResult struct {
	session.Session
}

func (SessionResult) IsResultObject() {}

// helper for unmarshaling dynamic xResults objects
type rawMap map[string]json.RawMessage

// implements ResultObject interface
// used for handling the raw RPC response
type AddResult struct {
	Section string `json:"section"`
}

func (AddResult) IsResultObject() {}

type Change []string
type ChangesResult struct {
	Many map[string][]Change `json:"many,omitempty"`
	One  []Change            `json:"one,omitempty"`
}

func (ChangesResult) IsResultObject() {}

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
func (v ChangesResult) MarshalJSON() ([]byte, error) {
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

func (v *ChangesResult) UnmarshalJSON(data []byte) (err error) {
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
type ConfigsResult struct {
	Configs []string `json:"configs"`
}

func (ConfigsResult) IsResultObject() {}

// implements ResultObject interface
// used for handling the raw RPC response
type ValueResult struct {
	Value uci.List `json:"value"`
}

func (ValueResult) IsResultObject() {}

type DataResult struct {
	Data string `json:"data"`
}

func (DataResult) IsResultObject() {}

type FileResult struct {
	Name  string `json:"name,omitempty"`
	Path  string `json:"path,omitempty"`
	Type  string `json:"type"`
	Size  int    `json:"size"`
	Mode  int    `json:"mode"`
	Atime int    `json:"atime"`
	Mtime int    `json:"mtime"`
	Ctime int    `json:"ctime"`
	Inode int    `json:"inode"`
	UID   int    `json:"uid"`
	GID   int    `json:"gid"`
}

func (FileResult) IsResultObject() {}

type EntriesResult struct {
	Entries []interface{} `json:"entries"`
}

func (EntriesResult) IsResultObject() {}

type MD5Result struct {
	MD5 string `json:"md5"`
}

func (MD5Result) IsResultObject() {}

type ExecResult struct {
	Code   int    `json:"code"`
	Stdout string `json:"stdout"`
}

func (ExecResult) IsResultObject() {}

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
type ValuesResult struct {
	Values map[string]uci.ConfigSection `json:"values"`
}

func (ValuesResult) IsResultObject() {}

func (v ValuesResult) MarshalJSON() ([]byte, error) {
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

func (v *ValuesResult) UnmarshalJSON(data []byte) (err error) {
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
	var val AddResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Section) > 0 {
			return val, nil
		}
	}

	return nil, nil
}

// checker for ExitCode
func matchExitCode(data json.RawMessage) (ResultObject, error) {
	var val ExitCode

	if err := json.Unmarshal(data, &val); err == nil {
		return val, nil
	}

	return nil, nil
}

// matcher for changesResult
func matchChangesResult(data json.RawMessage) (ResultObject, error) {
	var raw rawMap
	var val ChangesResult

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
	var val ConfigsResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Configs) > 0 {
			return val, nil
		}
	}

	return nil, nil
}

// checker for sessionResponse
func matchSessionResult(data json.RawMessage) (ResultObject, error) {
	var val session.Session

	if err := json.Unmarshal(data, &val); err == nil {
		if val.SessionID != "" { // easiest way to see if it unmarshaled into an empty Session struct
			return SessionResult{val}, nil
		}
	}
	return nil, nil
}

// matcher for valueResult
func matchValueResult(data json.RawMessage) (ResultObject, error) {
	var val ValueResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Value) > 0 {
			return val, nil
		}
	}

	return nil, nil
}

// matcher for valuesResult
func matchValuesResult(data json.RawMessage) (ResultObject, error) {
	var val ValuesResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Values) > 0 {
			return val, nil
		}
	}

	return nil, nil
}

func matchDataResult(data json.RawMessage) (ResultObject, error) {
	var val DataResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Data) > 0 {
			return val, nil
		}
	}

	return nil, nil
}

func matchEntriesResult(data json.RawMessage) (ResultObject, error) {
	var val EntriesResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Entries) > 0 {
			return val, nil
		}
	}

	return nil, nil
}

func matchFileResult(data json.RawMessage) (ResultObject, error) {
	var val FileResult

	if err := json.Unmarshal(data, &val); err == nil {
		if val != (FileResult{}) {
			return val, nil
		}
	}

	return nil, nil
}

func matchMD5Result(data json.RawMessage) (ResultObject, error) {
	var val MD5Result

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.MD5) > 0 {
			return val, nil
		}
	}

	return nil, nil
}

func matchExecResult(data json.RawMessage) (ResultObject, error) {
	var val ExecResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Stdout) > 0 {
			return val, nil
		}
	}

	return nil, nil
}

// response type registry
type resultObjectMatcher func(json.RawMessage) (ResultObject, error)

var resultObjectMatcherRegistry []resultObjectMatcher

func registerResultObjectMatcher(checker resultObjectMatcher) {
	resultObjectMatcherRegistry = append(resultObjectMatcherRegistry, checker)
}

// for all matchX funcs:
//
//	return (nil, nil) for non-matches
//	return (obj, nil) for valid matches
//	only return (nil, err) for broken JSON, which should almost never happen unless data is corrupted
func init() {
	registerResultObjectMatcher(matchExitCode)
	registerResultObjectMatcher(matchAddResult)
	registerResultObjectMatcher(matchChangesResult)
	registerResultObjectMatcher(matchConfigsResult)
	registerResultObjectMatcher(matchDataResult)
	registerResultObjectMatcher(matchEntriesResult)
	registerResultObjectMatcher(matchMD5Result)
	registerResultObjectMatcher(matchExecResult)
	registerResultObjectMatcher(matchFileResult)
	registerResultObjectMatcher(matchSessionResult)
	registerResultObjectMatcher(matchValueResult)
	registerResultObjectMatcher(matchValuesResult)
}
