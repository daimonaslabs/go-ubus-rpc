package uci

import (
	"encoding/json"
	"fmt"
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

type UCIConfigSectionOptions interface {
	IsUCIConfigSectionOptions()
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

type StringBool string

const (
	StringBoolFalse = StringBool("0")
	StringBoolTrue  = StringBool("1")
)

func ToStringBool(b bool) (s StringBool) {
	if b {
		return StringBoolTrue
	} else {
		return StringBoolFalse
	}
}
