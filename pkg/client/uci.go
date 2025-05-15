package client

import (
	"encoding/json"
	"errors"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
)

type UCIInterface interface {
	Configs() Call
	Get(opts *UCIGetOptions) Call
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
	c.setSignature(&UCIGetOptions{})

	return c.Call
}

func (c *uciCall) Get(opts *UCIGetOptions) Call {
	c.setProcedure("get")
	c.setSignature(opts)

	return c.Call
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
			return UCIConfigsResult{}, errors.New("not a ConfigsResult")
		}
	} else { // error
		return UCIConfigsResult{}, errors.New(p[0].(ExitCode).Error())
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
		case valueResult:
			u.Option = map[string]string{opts.Option: obj.Value}
		case valuesResult:
			err = json.Unmarshal(data, &u.SectionArray) // TODO need more custom unmarshalling since SectionArray is an interface
		default:
			return UCIGetResult{}, errors.New("not a GetResult")
		}
	} else { // error
		return UCIGetResult{}, errors.New(p[0].(ExitCode).Error())
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
type valueResult struct {
	Value string `json:"value"`
}

func (valueResult) isResultObject() {}

// implements ResultObject interface
// used for handling the raw RPC response
type valuesResult struct {
	Values map[string]map[string]any `json:"values"`
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

// matcher for valueResult
func matchValueResult(data json.RawMessage) (ResultObject, error) {
	var val valueResult

	if err := json.Unmarshal(data, &val); err == nil {
		if val.Value != "" {
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
			return valuesResult{Values: val.Values}, nil
		}
	}
	return nil, nil
}
