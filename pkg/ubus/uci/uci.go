package uci

import (
	"encoding/json"
	"fmt"
)

var (
	Configs = []string{
		"dhcp",
		"dropbear",
		"firewall",
		"luci",
		"network",
		"rpcd",
		"system",
		"ubootenv",
		"ucitrack",
		"uhttpd",
		"wireless"}
)

type UCIConfigSection interface {
	IsAnonymous() bool
	GetType() string
	GetName() string
	GetIndex() int
}

// implements UCIConfigSection
// implements json.Marshaler and json.Unmarshaler
type UCIConfigOptionsStatic struct {
	Anonymous bool   `json:".anonymous,omitempty"`
	Type      string `json:".type,omitempty"`
	Name      string `json:".name,omitempty"`
	Index     int    `json:".index,omitempty"`
}

func (s UCIConfigOptionsStatic) IsAnonymous() bool {
	return s.Anonymous
}

func (s UCIConfigOptionsStatic) GetType() string {
	return s.Type
}

func (s UCIConfigOptionsStatic) GetName() string {
	return s.Name
}

func (s UCIConfigOptionsStatic) GetIndex() int {
	return s.Index
}

// TODO figure out how to get a UCIConfigSection without the static options
// returns all the values of a UCIConfigSection with the static options stripped
func Values(s UCIConfigSection) (r UCIConfigSection) {
	//for option, value := range s {
	//	if option != ".anonymous" || option != ".type" || option != ".name" || option != ".index" {
	//		r[option] = value
	//	}
	//}

	sectionBytes, _ := json.Marshal(s)
	json.Unmarshal(sectionBytes, &r)

	return r
}

type DynamicList []string

// UnmarshalJSON allows dynamicList to accept either a single string or a list of strings.
func (d *DynamicList) UnmarshalJSON(data []byte) error {
	// try unmarshaling as a slice of strings first
	var list []string
	if err := json.Unmarshal(data, &list); err == nil {
		*d = DynamicList(list)
		return nil
	}

	// try unmarshaling as a single string
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*d = DynamicList{single}
		return nil
	}

	// if neither works, return an error
	return fmt.Errorf("dynamicList: unsupported JSON value: %s", string(data))
}

// MarshalJSON outputs as a string if there's one item, or a list if multiple
func (d DynamicList) MarshalJSON() ([]byte, error) {
	if len(d) == 1 {
		return json.Marshal(d[0])
	}
	return json.Marshal([]string(d))
}

var StringBoolFalse = "0"
var StringBoolTrue = "1"

// TODO make these safer, check for the right format within the strings. split into
// separate, more specific types as needed. see FirewallConfig.RuleSection for all the
// different time and date options needed.
type IP string
type MAC string
type StringBool string
type Time string

func ToStringBool(b bool) string {
	if b {
		return StringBoolTrue
	} else {
		return StringBoolFalse
	}
}

func True() string {
	return StringBoolTrue
}

func False() string {
	return StringBoolFalse
}
