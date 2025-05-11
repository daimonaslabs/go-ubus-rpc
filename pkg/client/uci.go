package client

import (
	"encoding/json"
)

func newUCICall(u *UbusRPC) *uciCall {
	u.uciCall.setSessionID(u.ubusSession.SessionID)
	u.uciCall.setPath("uci")
	return &u.uciCall
}

type UCIInterface interface {
	GetResult(r Response) UCIResult
	Configs() CallInterface
	Get(opts *UCIOptions) CallInterface
}

// implements UCIInterface
// implements CallInterface
type uciCall struct {
	Call
}

// TODO do all the type checking stuff in there to return a single one
func (c *uciCall) GetResult(r Response) UCIResult {
	return r[1].(UCIResult)
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
	Values map[string]any `json:"values"` // TODO possibly replace any with a more specific interface
}

func (UCIResult) isResultObject() {}

// implements ResultObject interface
// nested type used for JSON parsing
type valueResult struct {
	Value string `json:"value"`
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
