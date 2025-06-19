package uci

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	// the names of all the default configs managed by UCI
	DHCP     = "dhcp"
	Dropbear = "dropbear"
	Firewall = "firewall"
	LuCI     = "luci"
	Network  = "network"
	RPCD     = "rpcd"
	System   = "system"
	UBootEnv = "ubootenv"
	UCITrack = "ucitrack"
	UHTTPd   = "uhttpd"
	Wireless = "wireless"
)

var (
	Configs []string
)

func init() {
	Configs = []string{DHCP, Dropbear, Firewall, LuCI, Network, RPCD, System, UBootEnv, UCITrack, UHTTPd, Wireless}
}

type ConfigSection interface {
	IsAnonymous() bool
	GetType() string
	GetName() string
	GetIndex() int
}

// implements ConfigSection
// implements json.Marshaler and json.Unmarshaler
type StaticSectionOptions struct {
	Anonymous bool   `json:".anonymous,omitempty"`
	Type      string `json:".type,omitempty"`
	Name      string `json:".name,omitempty"`
	Index     int    `json:".index,omitempty"`
}

func (s StaticSectionOptions) IsAnonymous() bool {
	return s.Anonymous
}

func (s StaticSectionOptions) GetType() string {
	return s.Type
}

func (s StaticSectionOptions) GetName() string {
	return s.Name
}

func (s StaticSectionOptions) GetIndex() int {
	return s.Index
}

type ConfigSectionOptions interface {
	IsConfigSectionOptions()
}

// TODO change to Bool, make it a normal bool but un/marshals to/from strings
// also add a UCIInt, same thing but for ints
type Bool bool

// marshals the bool to a string value of "1" or "0"
func (b Bool) MarshalJSON() ([]byte, error) {
	if b {
		return []byte("\"1\""), nil
	} else {
		return []byte("\"0\""), nil
	}
}

// unmarshals from "1" or "0" back to true or false
func (b *Bool) UnmarshalJSON(data []byte) error {
	val := strings.ReplaceAll(string(data), "\"", "")
	if val == "1" {
		*b = Bool(true)
		fmt.Println(*b)
		return nil
	} else if val == "0" {
		*b = Bool(false)
		return nil
	} else {
		return errors.New("invalid string value for bool")
	}
}

type List []string

// MarshalJSON outputs as a string if there's one item, or a list if multiple
func (d List) MarshalJSON() ([]byte, error) {
	if len(d) == 1 {
		return json.Marshal(d[0])
	}
	return json.Marshal([]string(d))
}

// UnmarshalJSON allows List to accept either a single string or a list of strings.
func (d *List) UnmarshalJSON(data []byte) error {
	// try unmarshaling as a slice of strings first
	var list []string
	if err := json.Unmarshal(data, &list); err == nil {
		*d = List(list)
		return nil
	}

	// try unmarshaling as a single string
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*d = List{single}
		return nil
	}

	// if neither works, return an error
	return fmt.Errorf("List: unsupported JSON value: %s", string(data))
}
