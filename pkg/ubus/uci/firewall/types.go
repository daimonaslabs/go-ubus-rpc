package firewall

import (
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
)

// these are static values for the uci.UCIConfigOptionsStatic.Type field
const (
	DefaultsType   = "defaults"
	ForwardingType = "forwarding"
	RedirectType   = "redirect"
	RuleType       = "rule"
	ZoneType       = "zone"
)

// Used by RuleSection.ICMPType
var ICMPTypes = []string{
	"address-mask-reply",
	"address-mask-request",
	"any",
	"bad-header",
	"communication-prohibited",
	"destination-unreachable",
	"echo-reply",
	"echo-request",
	"fragmentation-needed",
	"host-precedence-violation",
	"host-prohibited",
	"host-redirect",
	"host-unknown",
	"host-unreachable",
	"ip-header-bad",
	"neighbour-advertisement",
	"network-prohibited",
	"network-redirect",
	"network-unknown",
	"network-unreachable",
	"packet-too-big",
	"parameter-problem",
	"ping",
	"pong",
	"port-unreachable",
	"precedence-cutoff",
	"protocol-unreachable",
	"redirect",
	"required-option-missing",
	"router-advertisement",
	"router-solicitation",
	"source-quench",
	"source-route-failed",
	"time-exceeded",
	"timestamp-reply",
	"timestamp-request",
	"TOS-host-redirect",
	"TOS-host-unreachable",
	"TOS-network-redirect",
	"TOS-network-unreachable",
	"ttl-exceeded",
	"ttl-zero-during-reassembly",
	"ttl-zero-during-transit",
	"unknown-header-type",
}

type DefaultsSection struct {
	uci.UCIConfigOptionsStatic `json:",inline"`
	// Accepts redirects. Implemented upstream in Linux Kernel.
	AcceptRedirects uci.StringBool `json:"accept_redirects,omitempty"`
	// Implemented upstream in Linux Kernel.
	AcceptSourceRoute uci.StringBool `json:"accept_source_route,omitempty"`
	// Determines method of packet rejection.
	AnyRejectCode int `json:"any_reject_code,omitempty"`
	// Enable Conntrack helpers.
	AutoHelper uci.StringBool `json:"auth_helper,omitempty"`
	// (fw4 only, OpenWRT 22.03 and later) Enable automatic nftables includes under /usr/share/nftables.d/
	AutoIncludes uci.StringBool `json:"auto_includes,omitempty"`
	// Enable generation of custom rule chain hooks for user generated rules. User rules would be typically
	// stored in firewall.user but some packages e.g. BCP38 also make use of these hooks.
	CustomChains uci.StringBool `json:"custom_chains,omitempty"`
	// Disable IPv6 firewall rules. (not supported by fw4).
	DisableIPv6 uci.StringBool `json:"disable_ipv6,omitempty"`
	// Drop invalid packets (e.g. not matching any active connection).
	DropInvalid uci.StringBool `json:"drop_invalid"`
	// Enable software flow offloading for connections (decrease cpu load / increase routing throughput).
	FlowOffloading uci.StringBool `json:"flow_offloading,omitempty"`
	// Enable hardware flow offloading for connecions (depends on flow_offloading and hw capability).
	FlowOffloadingHW uci.StringBool `json:"flow_offloading_hw,omitempty"`
	// Set policy for the FORWARD chain of the filter table.
	Forward string `json:"forward,omitempty"`
	// Set policy for the INPUT chain of the filter table.
	Input string `json:"input,omitempty"`
	// Set policy for the OUTPUT chain of the filter table.
	Output string `json:"output,omitempty"`
	// Enable SYN flood protection (obsoleted by synflood_protect setting).
	SynFlood uci.StringBool `json:"synFlood,omitempty" ubsu:"syn_flood,omitempty"`
	// Enable SYN flood protection.
	SynFloodProtect uci.StringBool `json:"synflood_protect,omitempty"`
	// Set rate limit (packets/second) for SYN packets above which the traffic is considered a flood.
	SynFloodRate string `json:"synflood_rate,omitempty"`
	// Set burst limit for SYN packets above which the traffic is considered a flood if it exceeds the allowed rate.
	SynFloodBurst string `json:"synflood_burst,omitempty"`
	// 0 Disable, 1 Enable, 2 Enable when requested for ingress (but disable for egress) Explicit Congestion
	// Notification. Affects only traffic originating from the router itself. Implemented upstream in Linux Kernel.
	TCPECN int `json:"tcp_ecn,omitempty"`
	// Enable the use of SYN cookies.
	TCPSynCookies uci.StringBool `json:"tcp_syncookies,omitempty"`
	// Determines method of packet rejection.
	TCPRejectCode int `json:"tcp_reject_code,omitempty"`
	// Enable TCP window scaling.
	TCPWindowScaling uci.StringBool `json:"tcp_window_scaling,omitempty"`
}

func (in *DefaultsSection) DeepCopyInto(out *DefaultsSection) {
	*out = *in
}

type ForwardingSection struct {
	uci.UCIConfigOptionsStatic `json:",inline"`
	// Specifies the traffic destination zone. Refers to one of the defined zone names.
	Dest string `json:"dest,omitempty"`
	// If set to 0, forward is disabled.
	Enabled uci.StringBool `json:"enabled,omitempty"`
	// Specifies the address family (ipv4, ipv6 or any) for which the rules are generated.
	Family string `json:"family,omitempty"`
	// If specified, match traffic against the given ipset. The match can be inverted by prefixing the value
	// with an exclamation mark.
	IPSet string `json:"ipset,omitempty"`
	// Unique forwarding name.
	Name string `json:"name,omitempty"`
	// Specifies the traffic source zone. Refers to one of the defined zone names. For typical port forwards this
	// usually is 'wan'.
	Src string `json:"src,omitempty"`
}

func (in *ForwardingSection) DeepCopyInto(out *ForwardingSection) {
	*out = *in
}

type RedirectSection struct {
	uci.UCIConfigOptionsStatic `json:",inline"`
	// Specifies the traffic destination zone. Refers to one of the defined zone names, or * for any zone. If
	// specified, the rule applies to forwarded traffic; otherwise, it is treated as input rule.
	Dest string `json:"dest,omitempty"`
	// Match incoming traffic directed to the specified destination IP address, CIDR notations can be used, see
	// note above. With no dest zone, this is treated as an input rule!
	DestIP string `json:"dest_ip,omitempty"`
	// Match incoming traffic directed at the given destination port or port range, if relevant proto is specified.
	// Multiple ports can be specified like '80 443 465' 1.
	DestPort string `json:"dest_port,omitempty"`
	// Enable the redirect rule or not.
	Enabled uci.StringBool `json:"enabled,omitempty"`
	// Specifies the address family (ipv4, ipv6 or any) for which the rules are generated. If unspecified, matches
	// the address family of other options in this section and defaults to ipv4.
	Family string `json:"family,omitempty"`
	Helper string `json:"helper,omitempty"`
	// If specified, match traffic against the given ipset. The match can be inverted by prefixing the value with an
	// exclamation mark. You can specify the direction as 'setname src' or 'setname dest'. The default if neither src
	// nor dest are added is to assume src.
	IPSet string `json:"ipset,omitempty"`
	// Maximum average matching rate; specified as a number, with an optional /second, /minute, /hour or /day suffix.
	// Examples: 3/second, 3/sec or 3/s.
	Limit string `json:"limit,omitempty"`
	// Maximum initial number of packets to match, allowing a short-term average above limit.
	LimitBurst int `json:"limit_burst,omitempty"`
	// If specified, match traffic against the given firewall mark, e.g. 0xFF to match mark 255 or 0x0/0x1 to match
	// any even mark value. The match can be inverted by prefixing the value with an exclamation mark, e.g. !0x10 to
	// match all but mark #16.
	Mark string `json:"mark,omitempty"`
	// If specified, only match traffic during the given days of the month, e.g. 2 5 30 to only match on every 2nd,
	// 5th and 30rd day of the month. The list can be inverted by prefixing it with an exclamation mark,
	// e.g. ! 31 to always match but on the 31st of the month.
	Monthdays uci.Time `json:"monthdays,omitempty"`
	// Name of redirect.
	Name string `json:"name,omitempty"`
	// Match incoming traffic using the given protocol. Can be one (or several when using list syntax) of tcp, udp,
	// udplite, icmp, esp, ah, sctp, or all or it can be a numeric value, representing one of these protocols or a
	// different one. A protocol name from /etc/protocols is also allowed. The number 0 is equivalent to all.
	Proto string `json:"proto,omitempty"`
	// Activate NAT reflection for this redirect - applicable to DNAT targets.
	Reflection uci.StringBool `json:"reflection,omitempty"`
	// The source address to use for NAT-reflected packets if reflection is 1. This can be internal or external,
	// specifying which interface’s address to use. Applicable to DNAT targets.
	ReflectionSrc string `json:"reflection_src,omitempty"`
	// List of zones for which reflection should be enabled. Applicable to DNAT targets.
	ReflectionZone []string `json:"reflection_zone,omitempty"`
	// Specifies the traffic source zone. Refers to one of the defined zone names. For typical port forwards this
	// usually is wan.
	Src string `json:"src,omitempty"`
	// For DNAT, match incoming traffic directed at the given destination IP address. For SNAT rewrite the source
	// address to the given address.
	SrcDIP uci.IP `json:"src_dip,omitempty"`
	// For DNAT, match incoming traffic directed at the given destination port or port range on this host. For
	// SNAT rewrite the source ports to the given value.
	SrcDPort string `json:"src_dport,omitempty"`
	// Match incoming traffic from the specified source IP address.
	SrcIP uci.IP `json:"src_ip,omitempty"`
	// Match incoming traffic from the specified MAC address.
	SrcMAC uci.MAC `json:"src_mac,omitempty"`
	// Match incoming traffic originating from the given source port or port range on the client host.
	SrcPort string `json:"src_port,omitempty"`
	// If specifed, only match traffic after the given date (inclusive).
	StartDate uci.Time `json:"start_date,omitempty"`
	// If specified, only match traffic after the given time of day (inclusive).
	StartTime uci.Time `json:"start_time,omitempty"`
	// If specified, only match traffic before the given date (inclusive).
	StopDate uci.Time `json:"stop_date,omitempty"`
	// If specified, only match traffic before the given time of day (inclusive).
	StopTime uci.Time `json:"stop_time,omitempty"`
	// If specified, only match traffic during the given week days, e.g. sun mon thu fri to only match on Sundays,
	// Mondays, Thursdays and Fridays. The list can be inverted by prefixing it with an exclamation mark,
	// e.g. ! sat sun to always match but on Saturdays and Sundays.
	Weekdays uci.Time `json:"weekdays,omitempty"`
	// Firewall action (ACCEPT, REJECT, DROP, MARK, NOTRACK) for matched traffic.
	Target string `json:"target,omitempty"`
	// Treat all given time values as UTC time instead of local time.
	UTCTime uci.StringBool `json:"utc_time,omitempty"`
}

func (in *RedirectSection) DeepCopyInto(out *RedirectSection) {
	*out = *in
}

type RuleSection struct {
	uci.UCIConfigOptionsStatic `json:",inline"`
	// Specifies the traffic destination zone. Refers to one of the defined zone names, or * for any zone. If
	// specified, the rule applies to forwarded traffic; otherwise, it is treated as input rule.
	Dest string `json:"dest,omitempty"`
	// Match incoming traffic directed to the specified destination IP address, CIDR notations can be used, see
	// note above. With no dest zone, this is treated as an input rule!
	DestIP string `json:"dest_ip,omitempty"`
	// Match incoming traffic directed at the given destination port or port range, if relevant proto is specified.
	// Multiple ports can be specified like '80 443 465' 1.
	DestPort  string `json:"dest_port,omitempty"`
	Device    string `json:"device,omitempty"`
	Direction string `json:"direction,omitempty"`
	// Enable or disable rule.
	Enabled uci.StringBool `json:"enabled,omitempty"`
	// Specifies the address family (ipv4, ipv6 or any) for which the rules are generated. If unspecified, matches
	// the address family of other options in this section and defaults to any.
	Family string `json:"family,omitempty"`
	Helper string `json:"helper,omitempty"`
	// For protocol icmp select specific ICMP types to match. Values can be either exact ICMP type numbers or type
	// names (see ICMPTypes var).
	ICMPType uci.DynamicList `json:"icmp_type,omitempty"`
	// If specified, match traffic against the given ipset. The match can be inverted by prefixing the value with an
	// exclamation mark. You can specify the direction as 'setname src' or 'setname dest'. The default if neither src
	// nor dest are added is to assume src.
	IPSet string `json:"ipset,omitempty"`
	// Maximum average matching rate; specified as a number, with an optional /second, /minute, /hour or /day suffix.
	// Examples: 3/minute, 3/min or 3/m.
	Limit string `json:"limit,omitempty"`
	// Maximum initial number of packets to match, allowing a short-term average above limit.
	LimitBurst int `json:"limit_burst,omitempty"`
	// If specified, match traffic against the given firewall mark, e.g. 0xFF to match mark 255 or 0x0/0x1 to match
	// any even mark value. The match can be inverted by prefixing the value with an exclamation mark, e.g. !0x10 to
	// match all but mark #16.
	Mark string `json:"mark,omitempty"`
	// If specified, only match traffic during the given days of the month, e.g. 2 5 30 to only match on every 2nd,
	// 5th and 30rd day of the month. The list can be inverted by prefixing it with an exclamation mark,
	// e.g. ! 31 to always match but on the 31st of the month.
	Monthdays uci.Time `json:"monthdays,omitempty"`
	// Name of rule.
	Name string `json:"name,omitempty"`
	// Match incoming traffic using the given protocol. Can be one (or several when using list syntax) of tcp,
	// udp, udplite, icmp, esp, ah, sctp, or all or it can be a numeric value, representing one of these protocols
	// or a different one. A protocol name from /etc/protocols is also allowed. The number 0 is equivalent to all.
	Proto string `json:"proto,omitempty"`
	// Zeroes out the bits given by mask and ORs value into the packet mark. If mask is omitted, 0xFFFFFFFF is
	// assumed.
	SetMark   string `json:"set_mark,omitempty"`
	SetHelper string `json:"set_helper,omitempty"`
	// Zeroes out the bits given by mask and XORs value into the packet mark. If mask is omitted, 0xFFFFFFFF is
	// assumed.
	SetXmark string `json:"set_xmark,omitempty"`
	// Specifies the traffic source zone. Refers to one of the defined zone names, or * for any zone. If omitted,
	// the rule applies to output traffic.
	Src string `json:"src,omitempty"`
	// Match incoming traffic from the specified source IP address, CIDR notations can be used, see note above.
	SrcIP uci.IP `json:"src_ip,omitempty"`
	// Match incoming traffic from the specified MAC address.
	SrcMAC uci.MAC `json:"src_mac,omitempty"`
	// Match incoming traffic from the specified source port or port range, if relevant proto is specified.
	// Multiple ports can be specified like '80 443 465' 1.
	SrcPort string `json:"src_port,omitempty"`
	// If specifed, only match traffic after the given date (inclusive).
	StartDate uci.Time `json:"start_date,omitempty"`
	// If specified, only match traffic after the given time of day (inclusive).
	StartTime uci.Time `json:"start_time,omitempty"`
	// If specified, only match traffic before the given date (inclusive).
	StopDate uci.Time `json:"stop_date,omitempty"`
	// If specified, only match traffic before the given time of day (inclusive).
	StopTime uci.Time `json:"stop_time,omitempty"`
	// Firewall action (ACCEPT, REJECT, DROP, MARK, NOTRACK) for matched traffic.
	Target string `json:"target,omitempty"`
	// Treat all given time values as UTC time instead of local time.
	UTCTime uci.StringBool `json:"utc_time,omitempty"`
	// If specified, only match traffic during the given week days, e.g. sun mon thu fri to only match on Sundays,
	// Mondays, Thursdays and Fridays. The list can be inverted by prefixing it with an exclamation mark,
	// e.g. ! sat sun to always match but on Saturdays and Sundays.
	Weekdays uci.Time `json:"weekdays,omitempty"`
}

func (in *RuleSection) DeepCopyInto(out *RuleSection) {
	*out = *in
}

type ZoneSection struct {
	uci.UCIConfigOptionsStatic `json:",inline"`
	// Add CT helpers for zone.
	AutoHelper uci.StringBool `json:"auto_helper,omitempty"`
	// Enable generation of custom rule chain hooks for user generated rules. Has no effect if disabled (0) in a
	// DefaultsSection.
	CustomChains uci.StringBool `json:"custom_chains,omitempty"`
	// List of L3 network interface names attached to this zone, e.g. tun+ or ppp+ to match any TUN or PPP interface.
	// This is specifically suitable for undeclared interfaces which lack built-in netifd support such as OpenVPN.
	// Otherwise network is preferable and device should be avoided.
	Device []string `json:"device,omitempty"`
	// If set to 0, zone is disabled
	Enabled uci.StringBool `json:"enabled,omitempty"`
	// Specifies the address family (ipv4, ipv6 or any) for which the rules are generated. If unspecified, matches
	// the address family of other options in this section and defaults to any.
	Family string `json:"family,omitempty"`
	// Policy (ACCEPT, REJECT, DROP) for forwarded zone traffic.
	Forward string `json:"forward,omitempty"`
	// List of helpers to add to zone.
	Helper []string `json:"helper,omitempty"`
	// Policy (ACCEPT, REJECT, DROP) for incoming zone traffic.
	Input string `json:"input,omitempty"`
	// Bit field to enable logging in the filter and/or mangle tables, bit 0 = filter, bit 1 = mangle.
	Log int `json:"log,omitempty"`
	// Limits the amount of log messages per interval.
	LogLimit string `json:"log_limit,omitempty"`
	// Specifies whether outgoing zone IPv4 traffic should be masqueraded. This is typically enabled on the wan zone.
	Masq uci.StringBool `json:"masq,omitempty"`
	// Specifies whether outgoing zone IPv6 traffic should be masqueraded. This is typically enabled on the wan zone.
	// Available with fw4. Requires sourcefilter=0 for DHCPv6 interfaces with missing GUA prefix.
	Masq6 uci.StringBool `json:"masq6,omitempty"`
	// Do not add DROP INVALID rules, if masquerading is used. The DROP rules are supposed to prevent NAT leakage.
	MasqAllowInvalid uci.StringBool `json:"masq_allow_invalid,omitempty"`
	// Limit masquerading to the given destination subnets. Negation is possible by prefixing the subnet with !;
	// multiple subnets are allowed.
	MasqDest []string `json:"masq_dest,omitempty"`
	// Limit masquerading to the given source subnets. Negation is possible by prefixing the subnet with !; multiple
	// subnets are allowed.
	MasqSrc []string `json:"masq_src,omitempty"`
	// Enable MSS clamping for outgoing zone traffic.
	MTUFix uci.StringBool `json:"mtu_fix,omitempty"`
	// Unique zone name. 11 characters is the maximum working firewall zone name length.
	Name string `json:"name,omitempty"`
	// List of interfaces attached to this zone. If omitted and neither extra* options, subnets nor devices are given,
	// the value of name is used by default. Alias interfaces defined in the network config cannot be used as valid
	// 'standalone' networks. Use list syntax.
	Network []string `json:"network,omitempty"`
	// Policy (ACCEPT, REJECT, DROP) for outgoing zone traffic.
	Output string `json:"output,omitempty"`
	// List of IP subnets attached to this zone.
	Subnet []string `json:"subnet,omitempty"`
}

func (in *ZoneSection) DeepCopyInto(out *ZoneSection) {
	*out = *in
}
