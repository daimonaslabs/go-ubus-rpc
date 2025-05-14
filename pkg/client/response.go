package client

import (
	"encoding/json"
	"fmt"
)

// interface for content within a Response
type ResultObject interface {
	isResultObject()
}

type resultStatic struct {
	Type   string         `json:"type"`
	Values map[string]any `json:"values"`
}

// effectively a tuple:
// Response[0] is always an int (ExitCode)
// Response[1] is always an xResult type (e.g. UCIResult)
type Response []ResultObject

// custom UnmarshalJSON for Response
func (r *Response) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, item := range raw {
		var matched bool
		for _, matcher := range resultObjectMatcherRegistry {
			if obj, err := matcher(item); err == nil && obj != nil {
				*r = append(*r, obj)
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf("unknown result object: %s", string(item))
		}
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

// implements ResultObject
// implements builtin.error
// always the first object of the Response tuple
type ExitCode int

func (ExitCode) isResultObject() {}

func (e ExitCode) Error() string {
	return fmt.Sprintf("exit status %d", e)
}

// checker for ExitCode
func matchExitCode(data json.RawMessage) (ResultObject, error) {
	var val int

	if err := json.Unmarshal(data, &val); err == nil {
		return ExitCode(val), nil
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
	registerResultObjectMatcher(matchSessionResult)
	registerResultObjectMatcher(matchValueResult)
	registerResultObjectMatcher(matchValuesResult)
	registerResultObjectMatcher(matchConfigsResult)
}
