package dhcp

import (
	"github.com/daimonaslabs/go-ubus-rpc/pkg/uci"
)

// TODO
type BootSection struct {
	uci.UCIConfigOptionsStatic `json:",inline"`
}

func (in *BootSection) DeepCopyInto(out *BootSection) {
	*out = *in
}

// TODO
type ClientSection struct {
	uci.UCIConfigOptionsStatic `json:",inline"`
}

func (in *ClientSection) DeepCopyInto(out *ClientSection) {
	*out = *in
}

// TODO
type DHCPSection struct {
	uci.UCIConfigOptionsStatic `json:",inline"`
}

func (in *DHCPSection) DeepCopyInto(out *DHCPSection) {
	*out = *in
}

type DnsmasqSection struct {
	uci.UCIConfigOptionsStatic `json:",inline"`
	// List of IP addresses for queried domains. See the dnsmasq man page for syntax details.
	Address []string `json:"address,omitempty"`
	// Add the local domain as search directive in resolv.conf.
	// +optional
	AddLocalDomain uci.StringBool `json:"add_local_domain,omitempty"`
	// Add A, AAAA, and PTR records for this router only on DHCP served LAN.
	// Enhanced function available since OpenWRT 18.06 with option AddLocalFQDN
	AddLocalHostname uci.StringBool `json:"add_local_hostname,omitempty"`
	// Add A, AAAA, and PTR records for this router only on DHCP served LAN.
	// 0: Disable.
	// 1: Hostname on Primary Address.
	// 2: Hostname on All Addresses.
	// 3: FDQN on All Addresses.
	// 4: iface.host.domain on All Addresses.
	// Available since OpenWRT 18.06
	AddLocalFQDN int `json:"add_local_fqdn,omitempty"`
	// Add the MAC address of the requester to DNS queries which are forwarded upstream; this may be used to do
	// DNS filtering by the upstream server.
	// The MAC address can only be added if the requester is on the same subnet as the dnsmasq server. Note that
	// the mechanism used to achieve this (an EDNS0 option) is not yet standardised, so this should be considered
	// experimental. Also note that exposing MAC addresses in this way may have security and privacy implications.
	// The string value must be either "base64" or "text".
	AddMAC string `json:"addmac,omitempty"`
	// Labels WAN interfaces like add_local_fqdn instead of your ISP assigned default which may be
	// obscure. WAN is inferred from config dhcp sections with option ignore 1 set, so they do not
	// need to be named WAN.
	// Available since OpenWRT 18.06
	AddWANFQDN int `json:"add_wan_fqdn,omitempty"`
	// Additional host files to read for serving DNS responses. Syntax in each file is the same as /etc/hosts.
	AddnHosts []string `json:"addnhosts,omitempty"`
	// Expose additional filesystem paths to the jailed dnsmasq process. This is useful in the case of manually
	// configured includes in the configuration file or symlinks pointing outside of the exposed paths as used,
	// for example, by an ad blocker or other name-banning package.
	AddnMount []string `json:"addnmount,omitempty"`
	// By default, when dnsmasq has more than one upstream server available, it will send queries to just one
	// server. Setting this parameter forces dnsmasq to send all queries to all available servers. The reply
	// from the server which answers first will be returned to the original requeser.
	AllServers uci.StringBool `json:"allservers,omitempty"`
	// Force dnsmasq into authoritative mode. This speeds up DHCP leasing. Used if this is the only server on
	// the network.
	Authoritative uci.StringBool `json:"authoritative,omitempty"`
	// IP addresses to convert into NXDOMAIN responses (to counteract “helpful” upstream DNS servers that never
	// return NXDOMAIN).
	BogusNXDOMAIN []uci.IP `json:"bogusnxdomain,omitempty"`
	// Reject reverse lookups to private IP ranges where no corresponding entry exists in /etc/hosts.
	BogusPriv uci.StringBool `json:"boguspriv,omitempty"`
	// When set to 0, use each network interface's DNS address in the local /etc/resolv.conf. Normally, only
	// the loopback address is used, and all queries go through dnsmasq.
	CacheLocal uci.StringBool `json:"cachelocal,omitempty"`
	// Size of dnsmasq query cache.
	CacheSize int `json:"cachesize,omitempty"`
	// Directory with additional configuration files.
	ConfDir string `json:"confdir,omitempty"`
	// Enable DBus messaging for dnsmasq.
	// Standard builds of dnsmasq on OpenWrt do not include DBus support.
	DBus uci.StringBool `json:"dbus,omitempty"`
	// Specifies BOOTP options, in most cases just the file name. You can also use:
	// “$FILENAME, $TFTP_SERVER_NAME, $TFTP_IP_ADDRESS”.
	DHCPBoot string `json:"dhcp_boot,omitempty"`
	// Maximum number of concurrent connections.
	DNSForwardMax int `json:"dnsforwardmax,omitempty"`
	// Specify an external file with per host DHCP options.
	DHCPHostsFile string `json:"dhcphostsfile,omitempty"`
	// Maximum number of DHCP leases.
	DHCPLeaseMax int `json:"dhcpleasemax,omitempty"`
	// Run a custom script upon DHCP lease add / renew / remove actions.
	DHCPScript string `json:"dhcpscript,omitempty"`
	// Validate DNS replies and cache DNSSEC data.
	// Requires the dnsmasq-full package. Please note that many applications now require DNSSEC to work properly,
	// e.g. Google apps on iOS like Gmail and Google Maps, and Windows Update and Windows Account activation on
	// Windows PCs.
	DNSSEC uci.StringBool `json:"dnssec,omitempty"`
	// Check the zones of unsigned replies to ensure that unsigned replies are allowed in those zones. This
	// protects against an attacker forging unsigned replies for signed DNS zones, but is slower and requires that
	// the nameservers upstream of dnsmasq are DNSSEC-capable.
	// Requires the dnsmasq-full package.
	// Caution: If you use this option on a device that doesn't have a hardware clock, DNS resolution may break
	// after a reboot of the device due to an incorrect system time.
	DNSSECCheckUnsigned uci.StringBool `json:"dnsseccheckunsigned,omitempty"`
	// DNS domain handed out to DHCP clients.
	Domain string `json:"domain,omitempty"`
	// Tells dnsmasq never to forward queries for plain names, without dots or domain parts, to upstream
	// nameservers. If the name is not known from /etc/hosts or DHCP then a “not found” answer is returned.
	DomainNeeded uci.StringBool `json:"domainneeded,omitempty"`
	// Specify the largest EDNS.0 UDP packet which is supported by the DNS forwarder.
	EDNSPacketMax int `json:"ednspacket_max,omitempty"`
	// Enable the builtin TFTP server.
	EnableTFTP uci.StringBool `json:"enable_tftp,omitempty"`
	// Add the local domain part to names found in /etc/hosts.
	ExpandHosts uci.StringBool `json:"expandhosts,omitempty"`
	// Do not forward requests that cannot be answered by public name servers.
	// Make sure it is disabled if you need to resolve SRV records or use SIP phones.
	FilterWin2k uci.StringBool `json:"filterwin2k,omitempty"`
	// Do not resolve unqualifed local hostnames. Needs Domain to be set.
	FQDN uci.StringBool `json:"fqdn,omitempty"`
	// List of interfaces to listen on. If unspecified, dnsmasq will listen to all interfaces except those listed
	// in NotInterface. Note that dnsmasq listens on loopback by default.
	Interface []string `json:"interface,omitempty"`
	// The syntax is: `list ipset '/example.com/example.org/example_ipv4,example_ipv6'``
	// Requires the dnsmasq-full package.
	IPSet []string `json:"ipset,omitempty"`
	// Store DHCP leases in this file.
	LeaseFile string `json:"leasefile,omitempty"`
	// Listen only on the specified IP addresses. If unspecified, listen on IP addresses from each interface.
	ListenAddress []uci.IP `json:"listen_address,omitempty"`
	// Look up DNS entries for this domain from /etc/hosts. This follows the same syntax as Server entries.
	// See the dnsmasq man page for more details.
	Local string `json:"local,omitempty"`
	// Choose IP address to match the incoming interface if multiple addresses are assigned to a host name in
	// /etc/hosts. Initially disabled, but still enabled in the config by default.
	LocaliseQueries uci.StringBool `json:"localise_queries,omitempty"`
	// Accept DNS queries only from hosts whose address is on a local subnet, ie a subnet for which an interface
	// exists on the server.
	LocalService uci.StringBool `json:"localservice,omitempty"`
	// Default TTL for locally authoritative answers.
	LocalTTL int `json:"local_ttl,omitempty"`
	// Use dnsmasq as a local system resolver. Depends on the NoResolv and ResolvFile options.
	LocalUse uci.StringBool `json:"localuse,omitempty"`
	// Enables extra DHCP logging; logs all the options sent to the DHCP clients and the tags used to determine
	// them.
	LogDHCP uci.StringBool `json:"logdhcp,omitempty"`
	// Set the facility to which dnsmasq will send syslog entries. See the dnsmasq man page for available
	// facilities.
	LogFacility string `json:"logfacility,omitempty"`
	// Log the results of DNS queries, dump cache on SIGUSR1, include requesting IP.
	LogQueries uci.StringBool `json:"logqueries,omitempty"`
	// Set the maximum TTL of DNS answers, even when the TTL in the answer is higher.
	MaxCacheTTL int `json:"max_cache_ttl,omitempty"`
	// Dnsmasq picks random ports as source for outbound queries. When this option is given, the ports used
	// will always be smaller than or equal to the specified MaxPort value (max valid value 65535). Useful for
	// systems behind firewalls.
	// See also MinPort.
	MaxPort int `json:"maxport,omitempty"`
	// Limit the TTL in the DNS answer to this value.
	MaxTTL int `json:"max_ttl,omitempty"`
	// Set the minimum TTL of DNS answers, even when the TTL in the answer is lower.
	MinCacheTTL int `json:"min_cache_ttl,omitempty"`
	// Dnsmasq picks random ports as source for outbound queries. When this option is given, the ports used
	// will always be larger than or equal to the specified MinPort value (min valid value 1024). Useful for
	// systems behind firewalls.
	// See also MaxPort.
	MinPort int `json:"minport,omitempty"`
	// Don't daemonize the dnsmasq process.
	NoDaemon uci.StringBool `json:"nodaemon,omitempty"`
	// Don't read DNS names from /etc/hosts.
	NoHosts uci.StringBool `json:"nohosts,omitempty"`
	// Disable caching of negative “no such domain” responses.
	NoNegCache uci.StringBool `json:"nonegcache,omitempty"`
	// By default dnsmasq checks if an IPv4 address is in use before allocating it to a host by sending ICMP
	// echo request (aka ping) to the address in question. This parameter allows to disable this check.
	NoPing uci.StringBool `json:"noping,omitempty"`
	// Don't read upstream servers from /etc/resolv.conf which is linked to resolvfile by default.
	NoResolv uci.StringBool `json:"noresolv,omitempty"`
	// Bind only configured interface addresses, instead of the wildcard address.
	NoWildcard uci.StringBool `json:"nowildcard,omitempty"`
	// Interfaces dnsmasq should not listen on.
	NotInterface []string `json:"notinterface,omitempty"`
	// Listening port for DNS queries, disables DNS server functionality if set to 0.
	Port int `json:"port,omitempty"`
	// Use a fixed port for outbound DNS queries.
	QueryPort int `json:"queryport,omitempty"`
	// Suppress logging of the routine operation of DHCP. Errors and problems will still be logged.
	QuietDHCP uci.StringBool `json:"quietdhcp,omitempty"`
	// Enable DHCPv4 Rapid Commit (fast address assignment) See RFC 4039.
	RapidCommit uci.StringBool `json:"rapidcommit,omitempty"`
	// Read static lease entries from /etc/ethers, re-read on SIGHUP.
	ReadEthers uci.StringBool `json:"readethers,omitempty"`
	// Enables DNS rebind attack protection by discarding upstream RFC1918 responses.
	RebindProtection uci.StringBool `json:"rebind_protection,omitempty"`
	// Allows upstream 127.0.0.0/8 responses, required for DNS based blacklist services, only takes effect if
	// rebind protection is enabled.
	RebindLocalhost uci.StringBool `json:"rebind_localhost,omitempty"`
	// List of domains to allow RFC1918 responses for, only takes effect if rebind protection is enabled.
	// The correct syntax is: `list rebind_domain '/example.com/'`
	RebindDomain []string `json:"rebind_domain,omitempty"`
	// Specifies an alternative resolv file.
	ResolvFile string `json:"resolvfile,omitempty"`
	// List of network range with a DNS server to forward reverse DNS requests to. See the dnsmasq man page
	// for syntax details.
	RevServer []string `json:"rev_server,omitempty"`
	// Dnsmasq is designed to choose IP addresses for DHCP clients using a hash of the client's MAC address.
	// This normally allows a client's address to remain stable long-term, even if the client sometimes allows
	// its DHCP lease to expire. In this default mode IP addresses are distributed pseudo-randomly over the
	// entire available address range. There are sometimes circumstances (typically server deployment) where
	// it is more convenient to have IP addresses allocated sequentially, starting from the lowest available
	// address, and setting this parameter enables this mode. Note that in the sequential mode, clients which
	// allow a lease to expire are much more likely to move IP address; for this reason it should not be
	// generally used.
	SequentialIP uci.StringBool `json:"sequential_ip,omitempty"`
	// List of DNS servers to forward requests to. See the dnsmasq man page for syntax details.
	Server []string `json:"server,omitempty"`
	// Specify upstream servers directly. If one or more optional domains are given, that server is used only
	// for those domains and they are queried only using the specified server.
	// Syntax is `server=/*.mydomain.tld/192.168.100.1` or see the dnsmasq man page for details.
	ServerList string `json:"serverlist,omitempty"`
	// Obey order of DNS servers in /etc/resolv.conf.
	StrictOrder uci.StringBool `json:"strictorder,omitempty"`
	// Specifies the TFTP root directory.
	TFTPRoot string `json:"tftp_root,omitempty"`
}

func (in *DnsmasqSection) DeepCopyInto(out *DnsmasqSection) {
	*out = *in
}

// TODO
type HostSection struct {
	uci.UCIConfigOptionsStatic `json:",inline"`
}

func (in *HostSection) DeepCopyInto(out *HostSection) {
	*out = *in
}

// TODO
type IPSetSection struct {
	uci.UCIConfigOptionsStatic `json:",inline"`
}

func (in *IPSetSection) DeepCopyInto(out *IPSetSection) {
	*out = *in
}

// TODO
type RelaySection struct {
	uci.UCIConfigOptionsStatic `json:",inline"`
}

func (in *RelaySection) DeepCopyInto(out *RelaySection) {
	*out = *in
}
