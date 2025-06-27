/*
Copyright 2025 Daimonas Labs.

Licensed under the GNU General Public License, Version 3 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.gnu.org/licenses/gpl-3.0.en.html

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package uci

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	// the names of all the default configs managed by UCI
	DHCP     = "dhcp"
	Dropbear = "dropbear"
	Firewall = "firewall"
	//LuCI     = "luci"
	//Network  = "network"
	//RPCD     = "rpcd"
	//System   = "system"
	//UBootEnv = "ubootenv"
	//UCITrack = "ucitrack"
	//UHTTPd   = "uhttpd"
	Wireless = "wireless"
)

var (
	Configs []string
)

func init() {
	Configs = []string{DHCP, Firewall, Wireless, Dropbear} //, LuCI, Network, RPCD, System, UBootEnv, UCITrack, UHTTPd,
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
	Anonymous bool   `json:".anonymous"`
	Type      string `json:".type"`
	Name      string `json:".name"`
	Index     int    `json:".index"`
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

type Bool bool

func BoolPtr(b bool) *Bool {
	ptr := Bool(b)
	return &ptr
}

// marshals the bool to a string value of "1" or "0"
func (b Bool) MarshalJSON() ([]byte, error) {
	if b {
		return []byte("\"1\""), nil
	} else {
		return []byte("\"0\""), nil
	}
}

// unmarshals from "1"/"on" or "0"/"off" back to true or false
func (b *Bool) UnmarshalJSON(data []byte) error {
	val := strings.ReplaceAll(string(data), "\"", "")
	switch val {
	case "1", "on":
		*b = Bool(true)
		return nil
	case "0", "off":
		*b = Bool(false)
		return nil
	default:
		return errors.New("invalid string value for bool")
	}
}

type Int int

func IntPtr(i int) *Int {
	ptr := Int(i)
	return &ptr
}

// marshals int to a string
func (i Int) MarshalJSON() ([]byte, error) {
	str := strconv.Itoa(int(i))

	return json.Marshal(str)
}

// unmarshals a string back to an int
func (i *Int) UnmarshalJSON(data []byte) (err error) {
	var str string
	var val int

	if err = json.Unmarshal(data, &str); err == nil {
		if val, err = strconv.Atoi(str); err == nil {
			*i = Int(val)
		}
	}

	return err
}

type List []string

// marshals as a string if there's one item or a list if multiple
func (d List) MarshalJSON() ([]byte, error) {
	if len(d) == 1 {
		return json.Marshal(d[0])
	}
	return json.Marshal([]string(d))
}

// accepts either a single string or a list of strings
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
