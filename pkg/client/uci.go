package client

import (
	"encoding/json"
	"errors"
)

const (
	ConfigsResultType = "configs"
	ValueResultType   = "value"
	ValuesResultType  = "values"
)

func newUCICall(u *UbusRPC) *uciCall {
	u.uciCall.setSessionID(u.ubusSession.SessionID)
	u.uciCall.setPath("uci")
	return &u.uciCall
}

type UCIInterface interface {
	GetResult(p Response) (u UCIResult, err error)
	Configs() CallInterface
	Get(opts *UCIOptions) CallInterface
}

// implements UCIInterface
// implements CallInterface
type uciCall struct {
	Call
}

// TODO could we really just have one Result type (resultStatic?)
// do all the type checking stuff in here to return a UCIResult type
func (c *uciCall) GetResult(p Response) (u UCIResult, err error) {
	if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case valueResult:
			u.Type = ValueResultType
			json.Unmarshal(data, &u.Values)
			return u, nil
		case valuesResult:
			u.Type = ValuesResultType
			json.Unmarshal(data, &u.Values)
			return u, nil
		case configsResult:
			u.Type = ConfigsResultType
			json.Unmarshal(data, &u.Values)
			return u, nil
		default:
			return UCIResult{}, errors.New("not a UCIResult")
		}
	} else { // error
		return UCIResult{}, errors.New(p[0].(ExitCode).Error())
	}
}

func (c *uciCall) Configs() CallInterface {
	c.setProcedure("configs")
	c.setSignature(&UCIOptions{})

	return c
}

func (c *uciCall) Get(opts *UCIOptions) CallInterface {
	c.setProcedure("get")
	c.setSignature(opts)

	return c
}

// implements Signature interface
type UCIOptions struct {
	Config  string `json:"config,omitempty"`
	Section string `json:"section,omitempty"`
	Type    string `json:"type,omitempty"`
	Option  string `json:"option,omitempty"`
}

func (UCIOptions) isOptsType() {}

// implements ResultObject interface
type UCIResult struct {
	resultStatic
}

//func (UCIResult) isResultObject() {}

// implements ResultObject interface
// nested type used for JSON parsing
type valueResult struct {
	Value string `json:"value"`
}

// implements ResultObject interface
// nested type used for JSON parsing
type configsResult struct {
	Configs []string `json:"configs"`
}

func (configsResult) isResultObject() {}

// matcher for valueResult
func matchConfigsResult(data json.RawMessage) (ResultObject, error) {
	var val configsResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Configs) > 0 {
			return configsResult{Configs: val.Configs}, nil
		}
	}
	return nil, nil
}

func (valueResult) isResultObject() {}

// matcher forvalueResult
func matchValueResult(data json.RawMessage) (ResultObject, error) {
	var val valueResult

	if err := json.Unmarshal(data, &val); err == nil {
		if val.Value != "" {
			return valueResult{Value: val.Value}, nil
		}
	}
	return nil, nil
}

// implements ResultObject interface
// nested type used for JSON parsing
type valuesResult struct {
	Values map[string]map[string]any `json:"values"`
}

func (valuesResult) isResultObject() {}

// matcher for valueResult
func matchValuesResult(data json.RawMessage) (ResultObject, error) {
	var val valuesResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Values) > 0 {
			return valuesResult{Values: val.Values}, nil
		}
	}
	return nil, nil
}
