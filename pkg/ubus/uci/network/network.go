package network

import (
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
)

const (
	// the name of this config
	Config = "network"

	// these are static values for the uci.StaticSectionOptions.Type field
	BridgeVLAN = "bridge-vlan"
	Device     = "device"
	Globals    = "globals"
	Interface  = "interface"
	Switch     = "switch"
	SwitchPort = "switch_port"
	SwitchVLAN = "switch_vlan"
)

var (
	Sections []string
)

func init() {
	Sections = []string{Device, Globals, Interface, Switch, SwitchVLAN}
}

// used by InterfaceSection.Proto
var Protocols = []string{
	"3g",
	"6in4",
	"6rd",
	"6to4",
	"aiccu",
	"dhcp",
	"dhcpv6",
	"dslite",
	"gre",
	"gretap",
	"grev6",
	"grev6tap",
	"hnet",
	"l2tp",
	"ncm",
	"none",
	"ppp",
	"pppoa",
	"pppoe",
	"pptp",
	"qmi",
	"relay",
	"static",
	"vti",
	"vtiv6",
	"vxlan",
	"wwan",
}

type BridgeVLANSection struct {
	uci.StaticSectionOptions `json:",inline"`
	BridgeVLANSectionOptions `json:",inline"`
}

func (in *BridgeVLANSection) DeepCopyInto(out *BridgeVLANSection) {
	*out = *in
}

type BridgeVLANSectionOptions struct {
	// The name of the device to associated the bridge VLAN with
	Device *string `json:"device,omitempty"`
	// A list of ports and whether tagged or untagged, each list entry must be of the
	// format 'lan1:t' or 'lan2:u*'
	Ports *uci.List `json:"ports,omitempty"`
	// The VLAN ID
	VLAN *uci.Int `json:"vlan,omitempty"`
}

func (BridgeVLANSectionOptions) IsConfigSectionOptions() {}

type DeviceSection struct {
	uci.StaticSectionOptions `json:",inline"`
	DeviceSectionOptions     `json:",inline"`
}

func (in *DeviceSection) DeepCopyInto(out *DeviceSection) {
	*out = *in
}

type DeviceSectionOptions struct {
	// Enables auto-negotiation of link parameters like speed and duplex.
	Autoneg *string `json:"autoneg,omitempty"`
	// Base L2 device (required for macvlan type).
	IfName *string `json:"ifname,omitempty"`
	// MAC address override for the device (e.g., 62:11:22:aa:bb:cc).
	MACAddr *string `json:"macaddr,omitempty"`
	// Logical name of the L3 device; must match the interface's device option.
	Name *string `json:"name,omitempty"`
	// List of L2 device names to be included in a bridge.
	Ports *uci.List `json:"ports,omitempty"`
	// Controls receive (RX) flow control. "1" enables RX pause frames.
	RxPause *string `json:"rxpause,omitempty"`
	// Set routing table name or number for type=vrf.
	Table *string `json:"table,omitempty"`
	// Device type, e.g., "bridge" or "vrf".
	Type *string `json:"type,omitempty"`
	// Controls transmission (TX) flow control. "1" enables TX pause frames.
	TxPause *string `json:"txpause,omitempty"`
}

func (DeviceSectionOptions) IsConfigSectionOptions() {}

type GlobalsSection struct {
	uci.StaticSectionOptions `json:",inline"`
	GlobalsSectionOptions    `json:",inline"`
}

func (in *GlobalsSection) DeepCopyInto(out *GlobalsSection) {
	*out = *in
}

type GlobalsSectionOptions struct {
	// Enables packet steering across CPUs:
	// 0 = disabled, 1 = enabled, 2 = enabled for all CPUs.
	PacketSteering *uci.Int `json:"packet_steering,omitempty"`
	// Toggles net.ipv4.tcp_l3mdev_accept (for VRF).
	TCPL3Mdev *uci.Bool `json:"tcp_l3mdev,omitempty"`
	// Toggles net.ipv4.udp_l3mdev_accept (for VRF).
	UDPL3Mdev *uci.Bool `json:"udp_l3mdev,omitempty"`
	// IPv6 ULA prefix; set to "auto" to generate automatically.
	ULAPrefix *string `json:"ula_prefix,omitempty"`
}

func (GlobalsSectionOptions) IsConfigSectionOptions() {}

type InterfaceSection struct {
	uci.StaticSectionOptions `json:",inline"`
	InterfaceSectionOptions  `json:",inline"`
}

func (in *InterfaceSection) DeepCopyInto(out *InterfaceSection) {
	*out = *in
}

type InterfaceSectionOptions struct {
	// Whether to bring up the interface on boot.
	Auto *uci.Bool `json:"auto,omitempty"`
	// Whether to disable the interface section entirely.
	Disabled *uci.Bool `json:"disabled,omitempty"`
	// Name of the associated L3 device (e.g., eth0.1, br-lan, tun0).
	// Must match the name in the corresponding device section.
	Device *string `json:"device,omitempty"`
	// Whether to assign IP settings even if the link is down.
	ForceLink *uci.Bool `json:"force_link,omitempty"`
	// Enables or disables IPv6 support on this interface.
	IPv6 *uci.Bool `json:"ipv6,omitempty"`
	// Name or number of the IPv4 routing table for this interface.
	IP4Table *string `json:"ip4table,omitempty"`
	// Name or number of the IPv6 routing table for this interface.
	IP6Table *string `json:"ip6table,omitempty"`
	// Override the default MTU for this interface.
	MTU *uci.Int `json:"mtu,omitempty"`

	//
	// bridge options
	//

	// Ageing time (in seconds) for dynamic MAC entries in the filtering database.
	AgeingTime *uci.Int `json:"ageing_time,omitempty"`
	// Whether to allow creating bridges with no ports.
	BridgeEmpty *bool `json:"bridge_empty,omitempty"`
	// Delay (in seconds) between port state transitions (STP); default 2, min 4 recommended.
	ForwardDelay *uci.Int `json:"forward_delay,omitempty"`
	// Size of the kernel multicast hash table.
	HashMax *uci.Int `json:"hash_max,omitempty"`
	// Interval (in 1/100s) between IGMP general queries.
	HelloTime *uci.Int `json:"hello_time,omitempty"`
	// Enables IGMP snooping to optimize multicast traffic distribution.
	IGMPSnooping *bool `json:"igmp_snooping,omitempty"`
	// Enables the bridge as an IGMP querier.
	MulticastQuerier *bool `json:"multicast_querier,omitempty"`
	// Maximum age (in seconds) before trying to become Root Bridge (STP).
	MaxAge *uci.Int `json:"max_age,omitempty"`
	// Priority value for the bridge (STP); lower is higher priority.
	Priority *uci.Int `json:"priority,omitempty"`
	// Interval (in 1/100s) for sending IGMP query responses after a leave group message.
	LastMemberInterval *uci.Int `json:"last_member_interval,omitempty"`
	// Interval (in 1/100s) between IGMP general queries.
	QueryInterval *uci.Int `json:"query_interval,omitempty"`
	// Interval (in 1/100s) within which IGMP query responses must be sent.
	QueryResponseInterval *uci.Int `json:"query_response_interval,omitempty"`
	// IGMP robustness value, influences query intervals and timeouts.
	Robustness *uci.Int `json:"robustness,omitempty"`
	// Enables the Spanning Tree Protocol (STP) to prevent network loops.
	STP *uci.Bool `json:"stp,omitempty"`
	// Enables VLAN-aware bridge mode.
	VLANFiltering *uci.Bool `json:"vlan_filtering,omitempty"`
}

func (InterfaceSectionOptions) IsConfigSectionOptions() {}

type SwitchSection struct {
	uci.StaticSectionOptions `json:",inline"`
	SwitchSectionOptions     `json:",inline"`
}

func (in *SwitchSection) DeepCopyInto(out *SwitchSection) {
	*out = *in
}

type SwitchSectionOptions struct {
	// Aging time (in seconds) for the ARL (MAC address) table. Default may differ by hardware.
	ARLAgeTime *uci.Int `json:"arl_age_time,omitempty"`
	// Enables VLAN-aware mode on the switch. Default may differ by hardware.
	EnableVLAN *uci.Bool `json:"enable_vlan,omitempty"`
	// Enables mirroring of received packets from source to monitor port.
	EnableMirrorRX *uci.Bool `json:"enable_mirror_rx,omitempty"`
	// Enables mirroring of transmitted packets from source to monitor port.
	EnableMirrorTX *uci.Bool `json:"enable_mirror_tx,omitempty"`
	// Enables IGMP snooping (behavior may vary).
	IGMPSnooping *uci.Bool `json:"igmp_snooping,omitempty"`
	// Enables IGMPv3 support (behavior may vary).
	IGMPv3 *uci.Bool `json:"igmp_v3,omitempty"`
	// Port to which mirrored packets are sent.
	MirrorMonitorPort *uci.Int `json:"mirror_monitor_port,omitempty"`
	// Port from which packets are mirrored.
	MirrorSourcePort *uci.Int `json:"mirror_source_port,omitempty"`
	// Name of the switch being configured.
	Name *string `json:"name,omitempty"`
	// Whether to reset the switch configuration.
	Reset *uci.Bool `json:"reset,omitempty"`
}

func (SwitchSectionOptions) IsConfigSectionOptions() {}

type SwitchPortSection struct {
	uci.StaticSectionOptions `json:",inline"`
	SwitchPortSectionOptions `json:",inline"`
}

func (in *SwitchPortSection) DeepCopyInto(out *SwitchPortSection) {
	*out = *in
}

type SwitchPortSectionOptions struct {
	// Enables Energy Efficient Ethernet (EEE) features to save power.
	EnableEEE *uci.Bool `json:"enable_eee,omitempty"`
	// Enables IGMP snooping on this port. Behavior may vary and is unconfirmed.
	IGMPSnooping *uci.Bool `json:"igmp_snooping,omitempty"`
	// Enables IGMPv3 on this port. Behavior may vary and is unconfirmed.
	IGMPv3 *uci.Bool `json:"igmp_v3,omitempty"`
	// Name of the switch device this port belongs to.
	Device *string `json:"device,omitempty"`
	// Index of the port to configure.
	Port *uci.Int `json:"port,omitempty"`
	// Port VLAN ID (PVID) to assign to untagged ingress packets.
	// This may refer to VLAN index or tag depending on platform behavior.
	PVID *uci.Int `json:"pvid,omitempty"`
}

func (SwitchPortSectionOptions) IsConfigSectionOptions() {}

type SwitchVLANSection struct {
	uci.StaticSectionOptions `json:",inline"`
	SwitchVLANSectionOptions `json:",inline"`
}

func (in *SwitchVLANSection) DeepCopyInto(out *SwitchVLANSection) {
	*out = *in
}

type SwitchVLANSectionOptions struct {
	// A human-readable description of the VLAN configuration.
	Description *string `json:"description,omitempty"`
	// The switch device to configure (must match a defined switch).
	Device *string `json:"device,omitempty"`
	// A string of space-separated port indices associated with the VLAN.
	// Use 't' suffix for tagged ports (e.g., "0 1 3t 5t").
	Ports *string `json:"ports,omitempty"`
	// VLAN tag number to use (VID); if unset, defaults to value of vlan.
	// VLANs 0 and 4095 may have special meaning.
	VID *uci.Int `json:"vid,omitempty"`
	// The VLAN table index to configure (not necessarily equal to VID).
	// May be limited depending on hardware
	VLAN *uci.Int `json:"vlan,omitempty"`
}

func (SwitchVLANSectionOptions) IsConfigSectionOptions() {}
