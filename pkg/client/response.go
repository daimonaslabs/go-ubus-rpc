package client

import (
	"encoding/json"
	"fmt"
)

// interface for content within a Response
type ResultObject interface {
	isResultObject()
}

// effectively a tuple:
// Response[0] is always an int (IntWrapper)
// Response[1] is always an xResult type (e.g. SessionResult)
type Response []ResultObject

// return the result of the response as the Go-type that it is
func GetAs[T ResultObject](r Response) (T, bool) {
	var zero T
	if len(r) < 2 {
		return zero, false
	}
	obj, ok := r[1].(T)
	return obj, ok
}

// custom UnmarshalJSON for Response
func (r *Response) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, item := range raw {
		var matched bool
		for _, checker := range resultTypeRegistry {
			if obj, err := checker(item); err == nil && obj != nil {
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

// used to wrap plain ints
type IntWrapper struct {
	Value int
}

func (IntWrapper) isResultObject() {}

// checker for IntWrapper
func matchIntWrapper(data json.RawMessage) (ResultObject, error) {
	var val int
	if err := json.Unmarshal(data, &val); err == nil {
		return IntWrapper{Value: val}, nil
	}
	//return nil, errors.New("not an IntWrapper")
	return nil, nil
}

// response type registry
type resultObjectChecker func(json.RawMessage) (ResultObject, error)

var resultTypeRegistry []resultObjectChecker

func registerResultType(checker resultObjectChecker) {
	resultTypeRegistry = append(resultTypeRegistry, checker)
}

// for all matchX funcs:
//
//	return (nil, nil) for non-matches
//	return (obj, nil) for valid matches
//	only return (nil, err) for broken JSON, which should almost never happen unless data is corrupted
func init() {
	registerResultType(matchIntWrapper)
	registerResultType(matchSessionResult)
	registerResultType(matchValueResult)
	registerResultType(matchValuesResult)
}
