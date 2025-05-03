package client

import (
	"encoding/json"
	"fmt"
)

type UCICallGetter interface {
	UCI() UCIInterface
}

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
	var tmp *ValueResult
	fmt.Println("ValueResult")

	if err := json.Unmarshal(data, &tmp); err != nil {
		return nil, err
	}
	if tmp.Value != "" {
		return &ValueResult{Value: tmp.Value}, nil
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
	var tmp *ValuesResult

	if err := json.Unmarshal(data, &tmp); err != nil {
		return nil, err
	}
	if len(tmp.Values) > 0 {
		return ValuesResult{Values: tmp.Values}, nil
	}
	return nil, nil
}
