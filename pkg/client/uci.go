package client

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/firewall"
)

type UCIInterface interface {
	Add(ctx context.Context, opts UCIAddOptions) (r Response, err error)
	Apply(ctx context.Context, opts UCIApplyOptions) (r Response, err error)
	//Changes(ctx context.Context, opts UCIChangesOptions) (r Response, err error)
	Configs(ctx context.Context, opts UCIConfigsOptions) (r Response, err error)
	//Delete(ctx context.Context) (r Response, err error)
	Get(ctx context.Context, opts UCIGetOptions) (r Response, err error)
	//Revert(ctx context.Context, opts UCIRevertOptions) (r Response, err error)
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

func (c *uciRPC) Configs(ctx context.Context, opts UCIConfigsOptions) (Response, error) {
	c.setProcedure("configs")
	c.setSignature(opts)

	return c.do(ctx)
}

func (c *uciRPC) Get(ctx context.Context, opts UCIGetOptions) (Response, error) {
	c.setProcedure("get")
	c.setSignature(opts)

	return c.do(ctx)
}

func (c *uciRPC) Set(ctx context.Context, opts UCISetOptions) (Response, error) {
	c.setProcedure("set")
	c.setSignature(opts)

	return c.do(ctx)
}

func (c *uciRPC) Apply(ctx context.Context, opts UCIApplyOptions) (Response, error) {
	c.setProcedure("apply")
	c.setSignature(opts)

	return c.do(ctx)
}

/*
################################################################
#
# all xOptions types are in this block. they all implement the
# Signature interface.
#
################################################################
*/

// implements Signature interface

type UCIAddOptions struct {
	Config uci.ConfigName  `json:"config,omitempty"`
	Type   uci.SectionType `json:"type,omitempty"`
}

func (UCIAddOptions) isOptsType() {}

func (opts UCIAddOptions) GetResult(p Response) (u UCIAddResult, err error) {
	if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case addResult:
			err = json.Unmarshal(data, &u)
		default:
			return u, errors.New("not an AddResult")
		}
	} else { // error
		return u, errors.New(p[0].(ExitCode).Error())
	}
	return u, err
}

// does not have a GetResult func because this command only returns the exit code
// implements Signature interface
type UCIApplyOptions struct {
	Rollback uci.StringBool `json:"rollback,omitempty"`
	Timeout  int            `json:"timeout,omitempty"`
}

func (UCIApplyOptions) isOptsType() {}

// implements Signature interface
// empty struct because there are no options but it has a special return type so we're
// following the same pattern as the other commands to get the result
type UCIConfigsOptions struct{}

func (UCIConfigsOptions) isOptsType() {}

func (opts UCIConfigsOptions) GetResult(p Response) (u UCIConfigsResult, err error) {
	if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case configsResult:
			err = json.Unmarshal(data, &u)
		default:
			return u, errors.New("not a ConfigsResult")
		}
	} else { // error
		return u, errors.New(p[0].(ExitCode).Error())
	}
	return u, err
}

// implements Signature interface
type UCIGetOptions struct {
	Config  uci.ConfigName `json:"config,omitempty"`
	Section string         `json:"section,omitempty"`
	Type    string         `json:"type,omitempty"`
	Option  string         `json:"option,omitempty"`
}

func (UCIGetOptions) isOptsType() {}

func (opts UCIGetOptions) GetResult(p Response) (u UCIGetResult, err error) {
	if len(p) > 1 {
		switch obj := p[1].(type) {
		case valueResult:
			u.Option = map[string]uci.DynamicList{opts.Option: obj.Value}
		case valuesResult:
			for _, section := range obj.Values {
				switch s := section.(type) {
				case firewall.DefaultsSection:
					u.SectionArray = append(u.SectionArray, s)
					s.Anonymous = false
				case firewall.ForwardingSection:
					u.SectionArray = append(u.SectionArray, s)
				case firewall.RedirectSection:
					u.SectionArray = append(u.SectionArray, s)
				case firewall.RuleSection:
					u.SectionArray = append(u.SectionArray, s)
				case firewall.ZoneSection:
					u.SectionArray = append(u.SectionArray, s)
				}
			}
		default:
			return u, errors.New("not a GetResult")
		}
	} else { // error
		return u, errors.New(p[0].(ExitCode).Error())
	}
	return u, err
}

// does not have a GetResult func because this command only returns the exit code
// implements Signature interface
type UCISetOptions struct {
	Config  uci.ConfigName              `json:"config,omitempty"`
	Section string                      `json:"section,omitempty"`
	Values  uci.UCIConfigSectionOptions `json:"values,omitempty"`
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

// result of a `uci configs` command
type UCIConfigsResult struct {
	Configs []string `json:"configs,omitempty"`
}

// result of a `uci get` command
type UCIGetResult struct {
	// if any combination of Config, Section, and Type are specified, return a set of
	// UCIConfigSection(s)
	SectionArray []uci.UCIConfigSection `json:"sectionArray,omitempty"`
	// if Option is specified, return a single option's value
	Option map[string]uci.DynamicList `json:"option,omitempty"`
}

/*
################################################################
#
# all unexported xResult types are in this block.
#
################################################################
*/

// implements ResultObject interface
// used for handling the raw RPC response
type addResult struct {
	Section string `json:"section"`
}

func (addResult) isResultObject() {}

// implements ResultObject interface
// used for handling the raw RPC response
type configsResult struct {
	Configs []string `json:"configs"`
}

func (configsResult) isResultObject() {}

// implements ResultObject interface
// used for handling the raw RPC response
type valueResult struct {
	Value uci.DynamicList `json:"value"`
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
// it will marshal it into the form of the first one but with only that object.
type valuesResult struct {
	Values map[string]uci.UCIConfigSection `json:"values"`
}

// helper for unmarshaling valuesResults
type rawMap map[string]json.RawMessage

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

	if isSingleObject(result) {
		section, err := unmarshalRawSection(values)
		if err != nil {
			return err
		}
		v.Values = map[string]uci.UCIConfigSection{section.GetName(): section}
		return nil
	} else {
		// handle named entries in map
		v.Values = make(map[string]uci.UCIConfigSection)
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

func unmarshalRawResult[S uci.UCIConfigSection](data []byte) (uci.UCIConfigSection, error) {
	var s S
	err := json.Unmarshal(data, &s)
	return s, err
}

func unmarshalRawSection(data []byte) (section uci.UCIConfigSection, err error) {
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
	case string(firewall.Defaults):
		section, err = unmarshalRawResult[firewall.DefaultsSection](data)
	case string(firewall.Forwarding):
		section, err = unmarshalRawResult[firewall.ForwardingSection](data)
	case string(firewall.Redirect):
		section, err = unmarshalRawResult[firewall.RedirectSection](data)
	case string(firewall.Rule):
		section, err = unmarshalRawResult[firewall.RuleSection](data)
	case string(firewall.Zone):
		section, err = unmarshalRawResult[firewall.ZoneSection](data)
	default:
		return nil, errors.New("invalid config section")
	}
	return section, err
}

// checks if the value of `values` is a single uci.UCIConfigSection or not
func isSingleObject(m map[string]json.RawMessage) bool {
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
