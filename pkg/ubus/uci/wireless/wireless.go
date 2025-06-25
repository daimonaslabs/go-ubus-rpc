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

package wireless

import (
	"encoding/json"
	"strconv"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
)

const (
	// the name of this config
	Config = "wireless"

	// these are static values for the uci.StaticSectionOptions.Type field
	WifiDevice  = "wifi-device"
	WifiIface   = "wifi-iface"
	WifiStation = "wifi-station"
	WifiVLAN    = "wifi-vlan"
)

var (
	Sections []string
)

func init() {
	Sections = []string{WifiDevice, WifiIface, WifiStation, WifiVLAN}
}

type Channel uint8

// marshals uint8 to a string
// if zero, marshals to "auto"
func (c Channel) MarshalJSON() ([]byte, error) {
	var str string
	if c == 0 {
		str = "auto"
	} else {
		str = strconv.Itoa(int(c))
	}

	return json.Marshal(str)
}

// unmarshals a string back to uint8
func (c *Channel) UnmarshalJSON(data []byte) (err error) {
	var str string
	var val int

	if err = json.Unmarshal(data, &str); err == nil {
		if val, err = strconv.Atoi(str); err == nil {
			if str == "auto" {
				*c = Channel(0)
			} else {
				*c = Channel(val)
			}
		}
	}

	return err
}

type WifiDeviceSection struct {
	uci.StaticSectionOptions `json:",inline"`
	WifiDeviceSectionOptions `json:",inline"`
}

func (in *WifiDeviceSection) DeepCopyInto(out *WifiDeviceSection) {
	*out = *in
}

type WifiDeviceSectionOptions struct {
	// Set the basic data rates. Each basic_rate is measured in kb/s. This option only has an effect on ap and adhoc wifi-ifaces.
	// It is recommended to use the cell_density option instead.
	BasicRate []uci.Int `json:"basic_rate,omitempty"`
	// Set the beacon interval in units of 1.024 ms. Valid range: 15–65535. Applies only to ap and adhoc wifi-ifaces.
	BeaconInt *uci.Int `json:"beacon_int,omitempty"`
	// Specifies the band: 2g, 5g, 6g, or 60g. Replaces hwmode (since 21.02.2).
	Band *string `json:"band,omitempty"`
	// Configures data rates based on coverage cell density. 0 = Disabled, 1 = Normal, 2 = High, 3 = Very High.
	CellDensity *uci.Int `json:"cell_density,omitempty"`
	// Specifies a narrow channel width in MHz, e.g. 5, 10, 20.
	ChanBW *uci.Int `json:"chanbw,omitempty"`
	// Specifies the wireless channel. “auto” means lowest available or ACS if supported.
	Channel *Channel `json:"channel,omitempty"` // string or int
	// Use specific channels when channel is in “auto” mode.
	Channels *uci.List `json:"channels,omitempty"`
	// Specifies the country code (e.g., "US", "DE"). Affects channel availability and power limits.
	Country *string `json:"country,omitempty"`
	// Enables 802.11d country IE advertisement. Requires country.
	CountryIE *uci.Bool `json:"country_ie,omitempty"`
	// Distance between the AP and the furthest client in meters.
	Distance *string `json:"distance,omitempty"`
	// Enables or disables the radio adapter. true = disabled.
	Disabled *uci.Bool `json:"disabled,omitempty"`
	// Enables automatic antenna selection by the driver.
	Diversity *uci.Bool `json:"diversity,omitempty"`
	// Specifies available HT/VHT/HE capabilities. Autodetected by driver.
	HTCapab *string `json:"ht_capab,omitempty"`
	// Specifies the high throughput mode: HT20, HT40, VHT80, HE20, etc.
	HTMode *string `json:"htmode,omitempty"`
	// Specifies the hardware mode. Deprecated in favor of Band since 21.02.2.
	HWMode *string `json:"hwmode,omitempty"`
	// Pass any custom options to hostapd-*.conf. Values are passed as-is.
	HostapdOptions *uci.List `json:"hostapd_options,omitempty"`
	// Specifies the MAC address of the radio adapter. Used to identify the interface.
	MACAddr *string `json:"macaddr,omitempty"`
	// Set the log level: 0=Verbose, 1=Debug, 2=Info, 3=Notice, 4=Warning.
	LogLevel *uci.Int `json:"log_level,omitempty"`
	// Allow legacy 802.11b data rates. 0 = Disallow, 1 = Allow.
	LegacyRates *uci.Bool `json:"legacy_rates,omitempty"`
	// Sets the minimum client capability mode required to connect. Options: n, ac.
	RequireMode *string `json:"require_mode,omitempty"`
	// Specifies the antenna for receiving. 0 = auto, 1 = antenna 1, 2 = antenna 2, etc.
	RXAntenna *uci.Int `json:"rxantenna,omitempty"`
	// Set the supported data rates in kb/s. Must be a superset of basic_rate.
	SupportedRates []uci.Int `json:"supported_rates,omitempty"`
	// Radio type: e.g., "mac80211", "broadcom". Usually autodetected.
	Type *string `json:"type,omitempty"`
	// Specifies the transmission power in dBm. Subject to regulatory limits.
	TxPower *uci.Int `json:"txpower,omitempty"`
	// Specifies the antenna for transmitting. Same values as rxantenna.
	TXAntenna *uci.Int `json:"txantenna,omitempty"`
	// Specifies the radio PHY device. Usually autodetected.
	PHY *string `json:"phy,omitempty"`

	//
	// MAC80211 options
	//

	// Reduction in antenna gain from regulatory maximum in dBi
	AntennaGain *uci.Int `json:"antenna_gain,omitempty"`
	// Fragmentation threshold
	Frag *uci.Int `json:"frag,omitempty"`
	// Disable honoring 40 MHz intolerance in coexistence flags of stations.
	// When enabled, the radio will continue using 40 MHz channels even if intolerance is indicated by another AP/station.
	HTCoex *uci.Int `json:"ht_coex,omitempty"`
	// Do not scan for overlapping BSSs in HT40+/- mode. May violate regulatory requirements if enabled.
	NoScan *uci.Bool `json:"noscan,omitempty"`
	// Alternative to phy used to identify the device based on /sys/devices path
	Path *string `json:"path,omitempty"`
	// Override the RTS/CTS threshold
	RTS *uci.Int `json:"rts,omitempty"`

	//
	// Broadcom options
	//

	// Enables Broadcom frame bursting (Xpress Technology) if supported
	FrameBurst *uci.Bool `json:"frameburst,omitempty"`
	// Limits the maximum allowed number of associated clients
	MaxAssoc *uci.Int `json:"maxassoc,omitempty"`
	// Slot time in milliseconds
	SlotTime *uci.Int `json:"slottime,omitempty"`

	//
	// Ubiquiti Nanostation options
	//

	// Specifies the antenna, possible values are 'vertical' for internal vertical polarization, 'horizontal' for internal
	// horizontal polarization, or 'external' to use the external antenna connector
	Antenna *string `json:"antenna,omitempty"`
}

func (WifiDeviceSectionOptions) IsConfigSectionOptions() {}

type WifiIfaceSection struct {
	uci.StaticSectionOptions `json:",inline"`
	WifiIfaceSectionOptions  `json:",inline"`
}

func (in *WifiIfaceSection) DeepCopyInto(out *WifiIfaceSection) {
	*out = *in
}

type WifiIfaceSectionOptions struct {
	// BSSID override (used in adhoc/sta/WDS mode).
	BSSID *string `json:"bssid,omitempty"`
	// List of supported basic data rates (in kb/s).
	BasicRate []uci.Int `json:"basic_rate,omitempty"`
	// Isolates clients across different radios on the same bridge.
	BridgeIsolate *uci.Bool `json:"bridge_isolate,omitempty"`
	// Delivery traffic indication message period (1–255).
	DTIMPeriod *uci.Int `json:"dtim_period,omitempty"`
	// Disables the wireless network interface.
	Disabled *uci.Bool `json:"disabled,omitempty"`
	// For STA: disables network block by default.
	DefaultDisabled *uci.Bool `json:"default_disabled,omitempty"`
	// Specifies the Wi-Fi adapter (must match a wifi-device).
	Device *string `json:"device"`
	// Enables 802.11h support.
	Doth *uci.Bool `json:"doth,omitempty"`
	// Sets WPA/WEP encryption method.
	Encryption *string `json:"encryption,omitempty"`
	// Custom interface name (max 15 characters).
	IfName *string `json:"ifname,omitempty"`
	// Enables client-to-client isolation on AP.
	Isolate *uci.Bool `json:"isolate,omitempty"`
	// Used for 802.11f IAPP.
	IAPPInterface *string `json:"iapp_interface,omitempty"`
	// Enables 802.11w (MFP): 0 = disabled, 1 = optional, 2 = required.
	IEEE80211w *uci.Int `json:"ieee80211w,omitempty"`
	// Maximum timeout for 802.11w SA Query (ms).
	IEEE80211wMaxTimeout *uci.Int `json:"ieee80211w_max_timeout,omitempty"`
	// Retry timeout for 802.11w SA Query (ms).
	IEEE80211wRetryTimeout *uci.Int `json:"ieee80211w_retry_timeout,omitempty"`
	// Additional hostapd BSS options.
	HostapdBSSOptions uci.List `json:"hostapd_bss_options,omitempty"`
	// Whether SSID should be hidden from beacon frames.
	Hidden *uci.Bool `json:"hidden,omitempty"`
	// Sets the maximum STA listen interval allowed.
	MaxListenInterval *uci.Int `json:"max_listen_int,omitempty"`
	// Maximum number of clients allowed to connect.
	MaxAssoc *uci.Int `json:"maxassoc,omitempty"`
	// Override MAC address, or set to "random".
	MACAddr *string `json:"macaddr,omitempty"`
	// MAC filter policy: disable, allow, deny.
	MACFilter *string `json:"macfilter,omitempty"`
	// List of MACs for the MAC filter.
	MACList *uci.List `json:"maclist,omitempty"`
	// Join this mesh ID (IEEE 802.11s).
	MeshID *string `json:"mesh_id,omitempty"`
	// Multicast rate (in kb/s). Only in mesh/adhoc.
	McastRate *uci.Int `json:"mcast_rate,omitempty"`
	// Operating mode: ap, sta, adhoc, mesh, monitor.
	Mode *string `json:"mode"`
	// Logical networks attached (L3 bridge).
	Network uci.List `json:"network"`
	// Opportunistic Wireless Encryption (OWE) BSSID.
	OWETransitionBSSID *string `json:"owe_transition_bssid,omitempty"`
	// Opportunistic Wireless Encryption (OWE) SSID.
	OWETransitionSSID *string `json:"owe_transition_ssid,omitempty"`
	// Operating Channel Validation (OCV) config: 0–2.
	OCV *uci.Int `json:"ocv,omitempty"`
	// Wireless passphrase or key (WPA/WEP).
	Key *string `json:"key,omitempty"`
	// WEP key slot 1.
	Key1 *string `json:"key1,omitempty"`
	// WEP key slot 2.
	Key2 *string `json:"key2,omitempty"`
	// WEP key slot 3.
	Key3 *string `json:"key3,omitempty"`
	// WEP key slot 4.
	Key4 *string `json:"key4,omitempty"`
	// RSN Preauthentication for WPA2-EAP.
	RSNPreauth *uci.Bool `json:"rsn_preauth,omitempty"`
	// Require MFP for SAE associations.
	SAERequireMFP *uci.Bool `json:"sae_require_mfp,omitempty"`
	// SAE PWE mechanism: 0 = hunting, 1 = hash, 2 = both.
	SAEPWE *uci.Int `json:"sae_pwe,omitempty"`
	// SSID to broadcast or connect to.
	SSID *string `json:"ssid"`
	// Start with AP beaconing disabled.
	StartDisabled *uci.Bool `json:"start_disabled,omitempty"`
	// Use of short preamble.
	ShortPreamble *uci.Bool `json:"short_preamble,omitempty"`
	// Supported data rates (in kb/s).
	SupportedRates []uci.Int `json:"supported_rates,omitempty"`
	// Enables 4-address mode (WDS).
	WDS *uci.Bool `json:"wds,omitempty"`
	// Enables Wi-Fi Multimedia (WMM) QoS mode.
	WMM *uci.Bool `json:"wmm,omitempty"`
}

func (WifiIfaceSectionOptions) IsConfigSectionOptions() {}
