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

package dropbear

import (
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
)

const (
	// the name of this config
	Config = "dropbear"

	// these are static values for the uci.StaticSectionOptions.Type field
	Dropbear = "dropbear"
)

var (
	Sections []string
)

func init() {
	Sections = []string{Dropbear}
}

type DropbearSection struct {
	uci.StaticSectionOptions `json:",inline"`
	DropbearSectionOptions   `json:",inline"`
}

func (in *DropbearSection) DeepCopyInto(out *DropbearSection) {
	*out = *in
}

type DropbearSectionOptions struct {
	// Path to a file to display before authentication.
	BannerFile *string `json:"BannerFile,omitempty"`
	// Set to 0 to disable password authentication.
	PasswordAuth *uci.Bool `json:"PasswordAuth,omitempty"`
	// Port to listen on for SSH connections.
	Port *uci.Int `json:"Port,omitempty"`
	// Set to 0 to disable root password authentication.
	RootPasswordAuth *uci.Bool `json:"RootPasswordAuth,omitempty"`
	// Set to 0 to disable SSH login as root.
	RootLogin *uci.Bool `json:"RootLogin,omitempty"`
	// Set to 1 to allow remote hosts to connect to forwarded ports.
	GatewayPorts *uci.Bool `json:"GatewayPorts,omitempty"`
	// Limit SSH access to a specific interface (e.g. "lan").
	Interface *string `json:"Interface,omitempty"`
	// Disconnect session after this many seconds of no activity (even with keepalives).
	IdleTimeout *uci.Int `json:"IdleTimeout,omitempty"`
	// List of key file paths for host keys.
	KeyFile *uci.List `json:"keyfile,omitempty"`
	// Set to 1 to announce the SSH service via mDNS.
	MDNS *uci.Int `json:"mdns,omitempty"`
	// Maximum allowed failed login attempts before connection closes.
	MaxAuthTries *uci.Int `json:"MaxAuthTries,omitempty"`
	// Set to 0 to disable starting Dropbear on system boot.
	Enable *uci.Bool `json:"enable,omitempty"`
	// Per-channel receive window buffer size.
	RecvWindowSize *uci.Int `json:"RecvWindowSize,omitempty"`
	// Deprecated: Path to RSA host key file (use keyfile instead).
	RSAKeyFile *string `json:"rsakeyfile,omitempty"`
	// Interval (in seconds) for SSH keepalives; 0 disables keepalives.
	SSHKeepAlive *uci.Int `json:"SSHKeepAlive,omitempty"`
	// Set to 1 to enable verbose output from the start script.
	Verbose *uci.Bool `json:"verbose,omitempty"`
}

func (DropbearSectionOptions) IsConfigSectionOptions() {}
