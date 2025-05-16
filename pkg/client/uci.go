package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/firewall"
)

type UCIInterface interface {
	Configs() Call
	Get(opts UCIGetOptions) Call
}

// implements UCIInterface
// implements CallInterface
type uciCall struct {
	Call
}

func newUCICall(u *UbusRPC) *uciCall {
	u.Call.setPath("uci")
	return &uciCall{u.Call}
}

func (c *uciCall) Configs() Call {
	c.setProcedure("configs")
	c.setSignature(UCIConfigsOptions{})

	return c.Call
}

func (c *uciCall) Get(opts UCIGetOptions) Call {
	c.setProcedure("get")
	c.setSignature(opts)

	return c.Call
}

func unmarshalRawResult[S uci.UCIConfigSection](data []byte) (uci.UCIConfigSection, error) {
	var s S
	err := json.Unmarshal(data, &s)
	return s, err
}

type rawSingleValues struct {
	Values json.RawMessage `json:"values"`
}

func unmarshalRawSection(data []byte) (section uci.UCIConfigSection, err error) {
	var probe struct {
		Type string `json:".type"`
	}
	var rawSV rawSingleValues

	if err = json.Unmarshal(data, &rawSV); err != nil {
		return nil, err
	}

	sectionBytes, err := json.Marshal(rawSV.Values)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(sectionBytes, &probe); err != nil {
		return nil, err
	}

	switch probe.Type {
	case firewall.DefaultsType:
		section, err = unmarshalRawResult[firewall.DefaultsSection](data)
	case firewall.ForwardingType:
		section, err = unmarshalRawResult[firewall.ForwardingSection](data)
	case firewall.RedirectType:
		section, err = unmarshalRawResult[firewall.RedirectSection](data)
	case firewall.RuleType:
		section, err = unmarshalRawResult[firewall.RuleSection](data)
	case firewall.ZoneType:
		section, err = unmarshalRawResult[firewall.ZoneSection](data)
	default:
		return nil, errors.New("invalid config section")
	}
	return section, err
}

type rawSectionArray struct {
	SectionArray map[string]json.RawMessage `json:"sections"`
}

func unmarshalRawSectionArrays(data []byte) (sections []uci.UCIConfigSection, err error) {
	//var probe struct {
	//	Type string `json:".type"`
	//}
	var rawMap rawSectionArray
	if err = json.Unmarshal(data, &rawMap); err != nil {
		return nil, err
	}

	for _, raw := range rawMap.SectionArray {
		section, err := unmarshalRawSection(raw)

		if err != nil {
			return nil, err
		}

		sections = append(sections, section)

		//if err = json.Unmarshal(raw, &probe); err != nil {
		//	return nil, err
		//}

		//switch probe.Type {
		//case firewall.DefaultsType:
		//	var s firewall.DefaultsSection
		//	if err = json.Unmarshal(raw, &s); err != nil {
		//		return nil, err
		//	}
		//	section = s
		//case firewall.ForwardingType:
		//	var s firewall.ForwardingSection
		//	if err = json.Unmarshal(raw, &s); err != nil {
		//		return nil, err
		//	}
		//	section = s
		//case firewall.RedirectType:
		//	var s firewall.RedirectSection
		//	if err = json.Unmarshal(raw, &s); err != nil {
		//		return nil, err
		//	}
		//	section = s
		//case firewall.RuleType:
		//	var s firewall.RuleSection
		//	if err = json.Unmarshal(raw, &s); err != nil {
		//		return nil, err
		//	}
		//	section = s
		//case firewall.ZoneType:
		//	var s firewall.ZoneSection
		//	if err = json.Unmarshal(raw, &s); err != nil {
		//		return nil, err
		//	}
		//	section = s
		//default:
		//	return nil, errors.New("invalid config section")
		//}
		//sections = append(sections, section)
	}
	return sections, nil
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
	Config  string `json:"config,omitempty"`
	Section string `json:"section,omitempty"`
	Type    string `json:"type,omitempty"`
	Option  string `json:"option,omitempty"`
}

func (UCIGetOptions) isOptsType() {}

func (opts UCIGetOptions) GetResult(p Response) (u UCIGetResult, err error) {
	if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch obj := p[1].(type) {
		case singleValuesResult:
			if opts.Config != "" && opts.Section != "" && opts.Option != "" {
				var value valueResult
				err = json.Unmarshal(data, &value)
				u.Option = map[string]string{opts.Option: value.Value}
			} else {
				section, err := unmarshalRawSection(data)
				if err != nil {
					return u, err
				}
				switch s := section.(type) {
				case firewall.DefaultsSection:
					u.SectionArray = map[string]uci.UCIConfigSection{section.GetName(): s}
				case firewall.ForwardingSection:
					u.SectionArray = map[string]uci.UCIConfigSection{section.GetName(): s}
				case firewall.RedirectSection:
					u.SectionArray = map[string]uci.UCIConfigSection{section.GetName(): s}
				case firewall.RuleSection:
					u.SectionArray = map[string]uci.UCIConfigSection{section.GetName(): s}
				case firewall.ZoneSection:
					u.SectionArray = map[string]uci.UCIConfigSection{section.GetName(): s}
				}
				fmt.Println("singleValuesResult: ", section, u)
			}
		case valueResult:
			u.Option = map[string]string{opts.Option: obj.Value}
		case valuesResult:
			sections, err := unmarshalRawSectionArrays(data)
			if err != nil {
				return u, err
			}
			for _, section := range sections {
				switch s := section.(type) {
				case firewall.DefaultsSection:
					u.SectionArray[s.Name] = s
				case firewall.ForwardingSection:
					u.SectionArray[s.Name] = s
				case firewall.RedirectSection:
					u.SectionArray[s.Name] = s
				case firewall.RuleSection:
					u.SectionArray[s.Name] = s
				case firewall.ZoneSection:
					u.SectionArray[s.Name] = s
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

/*
################################################################
#
# all xResult types are in this block.
#
################################################################
*/

// result of a `uci configs` command
type UCIConfigsResult struct {
	Configs []string `json:"configs,omitempty"`
}

// result of a `uci get` command
type UCIGetResult struct {
	SectionArray map[string]uci.UCIConfigSection `json:"sectionArray,omitempty"`
	Option       map[string]string               `json:"option,omitempty"`
}

// implements ResultObject interface
// used for handling the raw RPC response
type configsResult struct {
	Configs []string `json:"configs"`
}

func (configsResult) isResultObject() {}

// implements ResultObject interface
// used for handling the raw RPC response
type singleValuesResult struct {
	Values map[string]uci.UCIConfigSection `json:"values"`
}

func (singleValuesResult) isResultObject() {}

// implements ResultObject interface
// used for handling the raw RPC response
type valueResult struct {
	Value string `json:"value"`
}

func (valueResult) isResultObject() {}

// TODO make valuesResult implement json.(Un)Marshaler, call in
// unmarshalRawSection, change matchX funcs to check if the JSON
// has the correct top level key

// implements ResultObject interface
// used for handling the raw RPC response
type valuesResult struct {
	Values map[string]uci.UCIConfigSection `json:"values"`
	//Values map[string]map[string]any `json:"values"`
}

func (valuesResult) isResultObject() {}

/*
################################################################
#
# all matchXResult funcs are in this block. used in init().
#
################################################################
*/

// matcher for configsResult
func matchConfigsResult(data json.RawMessage) (ResultObject, error) {
	var val configsResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Configs) > 0 {
			return configsResult{Configs: val.Configs}, nil
		}
	}
	return nil, nil
}

// matcher for singleValuesResult
func matchSingleValuesResult(data json.RawMessage) (ResultObject, error) {
	var val singleValuesResult

	if err := json.Unmarshal(data, &val); err == nil {
		fmt.Println("matchSingleValuesResult: ", val)
		if val.Values != nil {
			return singleValuesResult{Values: val.Values}, nil
		}
	}
	return nil, nil
}

// matcher for valueResult
func matchValueResult(data json.RawMessage) (ResultObject, error) {
	var val valueResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Value) > 0 {
			return valueResult{Value: val.Value}, nil
		}
	}
	return nil, nil
}

// matcher for valuesResult
func matchValuesResult(data json.RawMessage) (ResultObject, error) {
	var val valuesResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Values) > 0 {
			fmt.Println("error????")
			return valuesResult{Values: val.Values}, nil
		}
	}
	return nil, nil
}
