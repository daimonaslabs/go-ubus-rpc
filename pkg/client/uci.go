package client

import (
	"encoding/json"
)

func newUCICall(u *UbusRPC) *uciCall {
	u.uciCall.SetSessionID(u.ubusSession.SessionID)
	u.uciCall.SetPath("uci")
	return &u.uciCall
}

type UCIInterface interface {
	Configs() CallInterface
	Get(opts *UCIOptions) CallInterface
}

// implements UCIInterface
type uciCall struct {
	Call
}

func (c *uciCall) Configs() CallInterface {
	c.SetProcedure("configs")
	c.SetSignature(map[string]any{})

	return c
}

func (c *uciCall) Get(opts *UCIOptions) CallInterface {
	c.SetProcedure("get")
	c.SetSignature(opts)

	return c
}

type UCIOptions struct {
	Config  string `json:"config,omitempty"`
	Section string `json:"section,omitempty"`
	Type    string `json:"type,omitempty"`
	Option  string `json:"option,omitempty"`
}

// implements ResultObject interface
type ValueResult struct {
	Value string `json:"value"`
}

func (ValueResult) isResultObject() {}

// checker for ValueResult
func matchValueResult(data json.RawMessage) (ResultObject, error) {
	var val ValueResult

	if err := json.Unmarshal(data, &val); err == nil {
		if val.Value != "" {
			return ValueResult{Value: val.Value}, nil
		}
	}
	return nil, nil
}

// implements ResultObject interface
type ValuesResult struct {
	Values map[string]map[string]any `json:"values"`
}

func (ValuesResult) isResultObject() {}

// checker for ValueResult
func matchValuesResult(data json.RawMessage) (ResultObject, error) {
	var val ValuesResult

	if err := json.Unmarshal(data, &val); err == nil {
		if len(val.Values) > 0 {
			return ValuesResult{Values: val.Values}, nil
		}
	}
	return nil, nil
}
