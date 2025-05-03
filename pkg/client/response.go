package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// interface for content within a Response
type ResultObject interface {
	isResultObject()
}

// effectively a tuple:
// Response[0] is always an int (IntWrapper)
// Response[1] is always an xResult type (e.g. SessionResult)
type Response []ResultObject

// custom UnmarshalJSON for Response
func (r *Response) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, item := range raw {
		var matched bool
		for _, checker := range resultTypeRegistry {
			if obj, err := checker(item); err == nil {
				*r = append(*r, obj)
				matched = true
				fmt.Println(matched, reflect.TypeOf(obj))
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
	return nil, errors.New("not an IntWrapper")
}

// response type registry
type resultObjectChecker func(json.RawMessage) (ResultObject, error)

var resultTypeRegistry []resultObjectChecker

func registerResultType(checker resultObjectChecker) {
	resultTypeRegistry = append(resultTypeRegistry, checker)
}

func init() {
	registerResultType(matchIntWrapper)
	registerResultType(matchSessionResult)
	registerResultType(matchValueResult)
	registerResultType(matchValuesResult)
}
