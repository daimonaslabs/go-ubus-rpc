package dhcp

import (
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
)

const (
	// the name of this config
	Config = "dhcp"

	// these are static values for the uci.StaticSectionOptions.Type field
	Boot        = "boot"
	CircuitID   = "circuitid"
	DHCP        = "dhcp"
	Dnsmasq     = "dnsmasq"
	Host        = "host"
	HostRecord  = "hostrecord"
	MAC         = "mac"
	Odhcpd      = "odhcpd"
	Relay       = "relay"
	RemoteID    = "remoteid"
	SubscrID    = "subscrid"
	Tag         = "tag"
	UserClass   = "userclass"
	VendorClass = "vendorclass"
)

var (
	Sections []string
)

func init() {
	Sections = []string{Boot, DHCP, Dnsmasq, Host, HostRecord, Relay}
}

type BootSection struct {
	uci.StaticSectionOptions `json:",inline"`
	BootSectionOptions       `json:",inline"`
}

func (in *BootSection) DeepCopyInto(out *BootSection) {
	*out = *in
}

type BootSectionOptions struct {
	// Additional options to be added for this network-id. If you specify this, you also need to specify
	// the network-id.
	DHCPOption uci.List `json:"dhcp_option,omitempty"`
	// The filename the host should request from the boot server.
	Filename string `json:"filename,omitempty"`
	// DHCPption will always be sent even if the client does not ask for it in the parameter request list. This
	// is sometimes needed, for example when sending options to PXELinux.
	Force uci.Bool `json:"force,omitempty"`
	// Dnsmasq instance to which the boot section is bound. If not specified the section is valid for all
	// dnsmasq instances.
	Instance string `json:"instance,omitempty"`
	// The tag (aka network-id) these boot options should apply to. Applies to all clients if left unspecified.
	NetworkID string `json:"networkid,omitempty"`
	// The IP address of the boot server.
	ServerAddress string `json:"serveraddress,omitempty"`
	// The hostname of the boot server.
	ServerName string `json:"servername,omitempty"`
}

func (BootSectionOptions) IsConfigSectionOptions() {}

type CircuitIDSection struct {
	uci.StaticSectionOptions `json:",inline"`
	CircuitIDSectionOptions  `json:",inline"`
}

func (in *CircuitIDSection) DeepCopyInto(out *CircuitIDSection) {
	*out = *in
}

type CircuitIDSectionOptions struct {
	// Matches the circuit ID as sent by the relay agent, as defined in RFC3046.
	CircuitID string `json:"circuitid,omitempty"`
	// The tag that matching clients will be assigned.
	NetworkID string `json:"networkid,omitempty"`
	// Additional options to be added for this tag aka networkid.
	DHCPOption uci.List `json:"dhcp_option,omitempty"`
	// Whether to send the additional options from dhcp_option list to the clients that didn't request them.
	Force uci.Bool `json:"force,omitempty"`
}

func (CircuitIDSectionOptions) IsConfigSectionOptions() {}

type DHCPSection struct {
	uci.StaticSectionOptions `json:",inline"`
	DHCPSectionOptions       `json:",inline"`
}

func (in *DHCPSection) DeepCopyInto(out *DHCPSection) {
	*out = *in
}

type DHCPSectionOptions struct {
	// Specifies whether DHCPv4 server should be enabled (server) or disabled (disabled).
	DHCPv4 string `json:"dhcpv4,omitempty"`
	// Specifies whether DHCPv6 server should be enabled (server), relayed (relay) or disabled (disabled).
	DHCPv6 string `json:"dhcpv6,omitempty"`
	// The ID dhcp_option here must be with written with an underscore. OpenWrt will translate this to
	// --dhcp-option, with a hyphen, as ultimately used by dnsmasq. Multiple option values can be given for
	// this network-id, with a a space between them and the total string between “”. E.g. '26,1470' or
	// 'option:mtu, 1470' that can assign an MTU per DHCP. Your client must accept MTU by DHCP for this to work.
	DHCPOption uci.List `json:"dhcp_option,omitempty"`
	// Exactly the same as dhcp_option (note the underscores), but it will be translated to --dhcp-option-force,
	// meaning that the DHCP option will be sent regardless on whether the client requested it.
	DHCPOptionForce []string `json:"dhcp_option_force,omitempty"`
	// DNS servers to announce on the network. Only IPv6 addresses are accepted. To configure IPv4 DNS servers,
	// use DHCPOption.
	DNS []string `json:"dns,omitempty"`
	// Announce the IPv6 address of interface as DNS service if the list of dns option is empty.
	DNSService uci.Bool `json:"dns_service,omitempty"`
	// Dynamically allocate client addresses, if set to 0 only clients present in the ethers files are served.
	DynamicDHCP uci.Bool `json:"dynamicdhcp,omitempty"`
	// Forces DHCP serving on the specified interface even if another DHCP server is detected on the same network
	// segment.
	Force uci.Bool `json:"force,omitempty"`
	// Specifies whether dnsmasq should ignore this pool if set to "1".
	Ignore uci.Bool `json:"ignore,omitempty"`
	// Dnsmasq instance to which the dhcp section is bound; if not specified the section is valid for all dnsmasq
	// instances.
	Instance string `json:"instance,omitempty"`
	// Specifies the interface associated with this DHCP address pool; must be one of the interfaces defined in
	// /etc/config/network.
	Interface string `json:"interface,omitempty"`
	// Specifies the lease time of addresses handed out to clients, for example 12h or 30m.
	LeaseTime string `json:"leasetime,omitempty"`
	// Specifies the size of the address pool (e.g. with start=100, limit=150, maximum address will be .249).
	Limit string `json:"limit,omitempty"`
	// Specifies whether DHCPv6, RA and NDP in relay mode is a master interface or not.
	Master uci.Bool `json:"master,omitempty"`
	// The dhcp functionality defined in the dhcp section is limited to the interface indicated here through
	// its network-id. In case omitted the system tries to know the network-id via the interface setting in this
	// dhcp section, through consultation of /etc/config/network. Some IDs get assigned dynamically, are not provided
	// by network, but still can be set here.
	NetworkID string `json:"networkid,omitempty"`
	// Specifies whether NDP should be relayed (relay) or disabled (disabled).
	NDP string `json:"ndp,omitempty"`
	// Ignore neighbor messages on slave enabled ("1") interfaces.
	NDProxySlave uci.Bool `json:"ndproxy_slave,omitempty"`
	// Learn routes from NDP.
	NDProxyRouting uci.Bool `json:"ndproxy_routing,omitempty"`
	// Specifies whether Router Advertisements should be enabled (server), relayed (relay) or disabled (disabled).
	RA string `json:"ra,omitempty"`
	// Default router lifetime in the RA message will be set if default route is present and a global IPv6 address (0)
	// or if default route is present but no global IPv6 address (1) or neither of both conditions (2).
	RADefault int `json:"ra_default,omitempty"`
	// List of RA flags to be advertised in RA messages:
	//	managed-config - get address and other information from DHCPv6 server. If this flag is set, other-config flag is redundant.
	//	other-config - get other configuration from DHCPv6 server (such as DNS servers).
	//	home-agent - see IETF docs for details.
	//	none.
	// OpenWrt since version 21.02 configures managed-config and other-config by default.
	RAFlags []string `json:"ra_flags,omitempty"`
	// Advertised current hop limit (0-255).
	RAHopLimit int `json:"ra_hoplimit,omitempty"`
	// Maximum advertised MTU.
	RAMTU int `json:"ra_mtu,omitempty"`
	// This option is deprecated. Use ra_flags and ra_slaac options instead.
	// RA management mode : no M-Flag but A-Flag (0), both M and A flags (1), M flag but not A flag (2).
	RAManagement int `json:"ra_management,omitempty"`
	// Maximum time interval between RAs (in seconds).
	RAMaxInterval int `json:"ra_maxinterval,omitempty"`
	// Minimum time interval between RAs (in seconds) .
	RAMinInterval int `json:"ra_mininterval,omitempty"`
	// Announce prefixes as offlink ("1") in RA messages.
	RAOfflink uci.Bool `json:"ra_offlink,omitempty"`
	// Announce routes with either high (high), medium (medium) or low (low) priority in RAs.
	RAPreference string `json:"ra_preference,omitempty"`
	// Advertised reachable time (in milliseconds) (0-3600000).
	RAReachableTime int `json:"ra_reachabletime,omitempty"`
	// Advertised NS retransmission time (in milliseconds) (0-60000).
	RARetransTime int `json:"ra_retranstime,omitempty"`
	// Announce DNS configuration in RA messages (RFC8106).
	RADNS uci.Bool `json:"ra_dns,omitempty"`
	// Announce SLAAC for a prefix (that is, set the A flag in RA messages).
	RASLAAC uci.Bool `json:"ra_slaac,omitempty"`
	// Advertised router lifetime (in seconds).
	RALifetime int `json:"ra_lifetime,omitempty"`
	// Limit the preferred and valid lifetimes of the prefixes in the RA messages to the configured DHCP leasetime.
	RAUseLeaseTime uci.Bool `json:"ra_useleasetime,omitempty"`
	// Specifies the offset from the network address of the underlying interface to calculate the minimum address that may be
	// leased to clients. It may be greater than 255 to span subnets.
	Start string `json:"start,omitempty"`
	// List of tags that dnsmasq needs to match to use with --dhcp-range.
	Tag []string `json:"tag,omitempty"`
}

func (DHCPSectionOptions) IsConfigSectionOptions() {}

type DnsmasqSection struct {
	uci.StaticSectionOptions `json:",inline"`
	DnsmasqSectionOptions    `json:",inline"`
}

func (in *DnsmasqSection) DeepCopyInto(out *DnsmasqSection) {
	*out = *in
}

type DnsmasqSectionOptions struct {
	uci.StaticSectionOptions `json:",inline"`
	// List of IP addresses for queried domains. See the dnsmasq man page for syntax details.
	Address []string `json:"address,omitempty"`
	// Add the local domain as search directive in resolv.conf.
	AddLocalDomain uci.Bool `json:"add_local_domain,omitempty"`
	// Add A, AAAA, and PTR records for this router only on DHCP served LAN.
	// Enhanced function available since OpenWRT 18.06 with option AddLocalFQDN
	AddLocalHostname uci.Bool `json:"add_local_hostname,omitempty"`
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
	AllServers uci.Bool `json:"allservers,omitempty"`
	// Force dnsmasq into authoritative mode. This speeds up DHCP leasing. Used if this is the only server on
	// the network.
	Authoritative uci.Bool `json:"authoritative,omitempty"`
	// IP addresses to convert into NXDOMAIN responses (to counteract “helpful” upstream DNS servers that never
	// return NXDOMAIN).
	BogusNXDOMAIN []string `json:"bogusnxdomain,omitempty"`
	// Reject reverse lookups to private IP ranges where no corresponding entry exists in /etc/hosts.
	BogusPriv uci.Bool `json:"boguspriv,omitempty"`
	// When set to 0, use each network interface's DNS address in the local /etc/resolv.conf. Normally, only
	// the loopback address is used, and all queries go through dnsmasq.
	CacheLocal uci.Bool `json:"cachelocal,omitempty"`
	// Size of dnsmasq query cache.
	CacheSize string `json:"cachesize,omitempty"`
	// Directory with additional configuration files.
	ConfDir string `json:"confdir,omitempty"`
	// Enable DBus messaging for dnsmasq.
	// Standard builds of dnsmasq on OpenWrt do not include DBus support.
	DBus uci.Bool `json:"dbus,omitempty"`
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
	// DNS domain handed out to DHCP clients.
	Domain string `json:"domain,omitempty"`
	// Tells dnsmasq never to forward queries for plain names, without dots or domain parts, to upstream
	// nameservers. If the name is not known from /etc/hosts or DHCP then a “not found” answer is returned.
	DomainNeeded uci.Bool `json:"domainneeded,omitempty"`
	// Specify the largest EDNS.0 UDP packet which is supported by the DNS forwarder.
	EDNSPacketMax string `json:"ednspacket_max,omitempty"`
	// Enable the builtin TFTP server.
	EnableTFTP uci.Bool `json:"enable_tftp,omitempty"`
	// Add the local domain part to names found in /etc/hosts.
	ExpandHosts uci.Bool `json:"expandhosts,omitempty"`
	// Do not forward requests that cannot be answered by public name servers.
	// Make sure it is disabled if you need to resolve SRV records or use SIP phones.
	FilterWin2k uci.Bool `json:"filterwin2k,omitempty"`
	// Do not resolve unqualifed local hostnames. Needs Domain to be set.
	FQDN uci.Bool `json:"fqdn,omitempty"`
	// List of interfaces to listen on. If unspecified, dnsmasq will listen to all interfaces except those listed
	// in NotInterface. Note that dnsmasq listens on loopback by default.
	Interface []string `json:"interface,omitempty"`
	// Store DHCP leases in this file.
	LeaseFile string `json:"leasefile,omitempty"`
	// Listen only on the specified IP addresses. If unspecified, listen on IP addresses from each interface.
	ListenAddress []string `json:"listen_address,omitempty"`
	// Look up DNS entries for this domain from /etc/hosts. This follows the same syntax as Server entries.
	// See the dnsmasq man page for more details.
	Local string `json:"local,omitempty"`
	// Choose IP address to match the incoming interface if multiple addresses are assigned to a host name in
	// /etc/hosts. Initially disabled, but still enabled in the config by default.
	LocaliseQueries uci.Bool `json:"localise_queries,omitempty"`
	// Accept DNS queries only from hosts whose address is on a local subnet, ie a subnet for which an interface
	// exists on the server.
	LocalService uci.Bool `json:"localservice,omitempty"`
	// Default TTL for locally authoritative answers.
	LocalTTL int `json:"local_ttl,omitempty"`
	// Use dnsmasq as a local system resolver. Depends on the NoResolv and ResolvFile options.
	LocalUse uci.Bool `json:"localuse,omitempty"`
	// Enables extra DHCP logging; logs all the options sent to the DHCP clients and the tags used to determine
	// them.
	LogDHCP uci.Bool `json:"logdhcp,omitempty"`
	// Set the facility to which dnsmasq will send syslog entries. See the dnsmasq man page for available
	// facilities.
	LogFacility string `json:"logfacility,omitempty"`
	// Log the results of DNS queries, dump cache on SIGUSR1, include requesting IP.
	LogQueries uci.Bool `json:"logqueries,omitempty"`
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
	NoDaemon uci.Bool `json:"nodaemon,omitempty"`
	// Don't read DNS names from /etc/hosts.
	NoHosts uci.Bool `json:"nohosts,omitempty"`
	// Disable caching of negative “no such domain” responses.
	NoNegCache uci.Bool `json:"nonegcache,omitempty"`
	// By default dnsmasq checks if an IPv4 address is in use before allocating it to a host by sending ICMP
	// echo request (aka ping) to the address in question. This parameter allows to disable this check.
	NoPing uci.Bool `json:"noping,omitempty"`
	// Don't read upstream servers from /etc/resolv.conf which is linked to resolvfile by default.
	NoResolv uci.Bool `json:"noresolv,omitempty"`
	// Bind only configured interface addresses, instead of the wildcard address.
	NonWildcard uci.Bool `json:"nonwildcard,omitempty"`
	// Interfaces dnsmasq should not listen on.
	NotInterface []string `json:"notinterface,omitempty"`
	// Listening port for DNS queries, disables DNS server functionality if set to 0.
	Port int `json:"port,omitempty"`
	// Use a fixed port for outbound DNS queries.
	QueryPort int `json:"queryport,omitempty"`
	// Suppress logging of the routine operation of DHCP. Errors and problems will still be logged.
	QuietDHCP uci.Bool `json:"quietdhcp,omitempty"`
	// Enable DHCPv4 Rapid Commit (fast address assignment) See RFC 4039.
	RapidCommit uci.Bool `json:"rapidcommit,omitempty"`
	// Read static lease entries from /etc/ethers, re-read on SIGHUP.
	ReadEthers uci.Bool `json:"readethers,omitempty"`
	// Enables DNS rebind attack protection by discarding upstream RFC1918 responses.
	RebindProtection uci.Bool `json:"rebind_protection,omitempty"`
	// Allows upstream 127.0.0.0/8 responses, required for DNS based blacklist services, only takes effect if
	// rebind protection is enabled.
	RebindLocalhost uci.Bool `json:"rebind_localhost,omitempty"`
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
	SequentialIP uci.Bool `json:"sequential_ip,omitempty"`
	// List of DNS servers to forward requests to. See the dnsmasq man page for syntax details.
	Server []string `json:"server,omitempty"`
	// Specify upstream servers directly. If one or more optional domains are given, that server is used only
	// for those domains and they are queried only using the specified server.
	// Syntax is `server=/*.mydomain.tld/192.168.100.1` or see the dnsmasq man page for details.
	ServerList string `json:"serverlist,omitempty"`
	// Obey order of DNS servers in /etc/resolv.conf.
	StrictOrder uci.Bool `json:"strictorder,omitempty"`
	// Specifies the TFTP root directory.
	TFTPRoot string `json:"tftp_root,omitempty"`
}

func (DnsmasqSectionOptions) IsConfigSectionOptions() {}

type HostSection struct {
	uci.StaticSectionOptions `json:",inline"`
	HostSectionOptions       `json:",inline"`
}

func (in *HostSection) DeepCopyInto(out *HostSection) {
	*out = *in
}

type HostSectionOptions struct {
	// Force broadcast DHCP response.
	Broadcast uci.Bool `json:"broadcast,omitempty"`
	// Add static forward and reverse DNS entries for this host.
	DNS uci.Bool `json:"dns,omitempty"`
	// The DHCPv6-DUID of this host.
	DUID string `json:"duid,omitempty"`
	// The IPv6 interface identifier (address suffix) as hexadecimal number (max. 16 chars, 64 bits, 8 bytes).
	HostID string `json:"hostid,omitempty"`
	// The IP address to be used for this host, or ignore to ignore any DHCP request from this host.
	IP string `json:"ip,omitempty"`
	// Dnsmasq instance to which the host section is bound; if not specified the section is valid for all dnsmasq instances.
	Instance string `json:"instance,omitempty"`
	// Host-specific lease time, e.g. 2m, 3h, 5d.
	LeaseTime string `json:"leasetime,omitempty"`
	// The hardware address(es) of this host, separated by spaces.
	MAC string `json:"mac,omitempty"`
	// If specified the section will apply only to requests having all the tags;
	// incoming interface name is always auto-assigned, other tags can be added by vendorclass/userclass/etc. sections.
	MatchTag []string `json:"match_tag,omitempty"`
	// Optional hostname to assign.
	Name string `json:"name,omitempty"`
	// Set the given tag for matching hosts.
	Tag string `json:"tag,omitempty"`
}

func (HostSectionOptions) IsConfigSectionOptions() {}

type HostRecordSection struct {
	uci.StaticSectionOptions `json:",inline"`
	HostRecordSectionOptions `json:",inline"`
}

func (in *HostRecordSection) DeepCopyInto(out *HostRecordSection) {
	*out = *in
}

type HostRecordSectionOptions struct {
	// The domain name.
	Name string `json:"name,omitempty"`
	// The IP address to resolve the name to.
	IP string `json:"ip,omitempty"`
}

func (HostRecordSectionOptions) IsConfigSectionOptions() {}

type MACSection struct {
	uci.StaticSectionOptions `json:",inline"`
	MACSectionOptions        `json:",inline"`
}

func (in *MACSection) DeepCopyInto(out *MACSection) {
	*out = *in
}

type MACSectionOptions struct {
	// Hardware address of the client.
	MAC string `json:"mac,omitempty"`
	// The tag that matching clients will be assigned.
	NetworkID string `json:"networkid,omitempty"`
	// Additional options to be added for this tag aka networkid.
	DHCPOption []string `json:"dhcp_option,omitempty"`
	// Whether to send the additional options from dhcp_option list to the clients that didn't request them.
	Force uci.Bool `json:"force,omitempty"`
}

func (MACSectionOptions) IsConfigSectionOptions() {}

type OdhcpdSection struct {
	uci.StaticSectionOptions `json:",inline"`
	OdhcpdSectionOptions     `json:",inline"`
}

func (in *OdhcpdSection) DeepCopyInto(out *OdhcpdSection) {
	*out = *in
}

type OdhcpdSectionOptions struct {
	// Use odhcpd as the main DHCPv4 service.
	MainDHCP uci.Bool `json:"maindhcp,omitempty"`
	// Location of the lease/hostfile for DHCPv4 and DHCPv6.
	LeaseFile string `json:"leasefile,omitempty"`
	// Location of the lease trigger script.
	LeaseTrigger string `json:"leasetrigger,omitempty"`
	// Enable DHCPv4 if the 'dhcp' section contains a start option, but no dhcpv4 option set.
	Legacy uci.Bool `json:"legacy,omitempty"`
	// Syslog level priority (0-7):
	// 0=emer, 1=alert, 2=crit, 3=err, 4=warn, 5=notice, 6=info, 7=debug
	LogLevel string `json:"loglevel,omitempty"` // TODO add a uci.UCIInt type and use here
}

func (OdhcpdSectionOptions) IsConfigSectionOptions() {}

type RelaySection struct {
	uci.StaticSectionOptions `json:",inline"`
	RelaySectionOptions      `json:",inline"`
}

func (in *RelaySection) DeepCopyInto(out *RelaySection) {
	*out = *in
}

type RelaySectionOptions struct {
	// A unique name for the section, which must be different to every other section's name.
	ID string `json:"id,omitempty"`
	// Logical network interface where the destination DHCP server is located.
	Interface string `json:"interface,omitempty"`
	// IP address to listen for DHCP requests.
	LocalAddr string `json:"local_addr,omitempty"`
	// IP address of the upstream DHCP server accessible through the network given by the interface option. DHCP
	// responses picked up on the far subnet will be relayed to this server. This address must be routed correctly
	// (i.e. you can ping it successfully from the OpenWrt command line).
	ServerAddr string `json:"server_addr,omitempty"`
}

func (RelaySectionOptions) IsConfigSectionOptions() {}

type RemoteIDSection struct {
	uci.StaticSectionOptions `json:",inline"`
	RemoteIDSectionOptions   `json:",inline"`
}

func (in *RemoteIDSection) DeepCopyInto(out *RemoteIDSection) {
	*out = *in
}

type RemoteIDSectionOptions struct {
	// Matches the remote ID as sent by the relay agent, as defined in RFC3046.
	RemoteID string `json:"remoteid,omitempty"`
	// The tag that matching clients will be assigned.
	NetworkID string `json:"networkid,omitempty"`
	// Additional options to be added for this tag aka networkid.
	DHCPOption []string `json:"dhcp_option,omitempty"`
	// Whether to send the additional options from dhcp_option list to the clients that didn't request them.
	Force uci.Bool `json:"force,omitempty"`
}

func (RemoteIDSectionOptions) IsConfigSectionOptions() {}

type SubscrIDSection struct {
	uci.StaticSectionOptions `json:",inline"`
	SubscrIDSectionOptions   `json:",inline"`
}

func (in *SubscrIDSection) DeepCopyInto(out *SubscrIDSection) {
	*out = *in
}

type SubscrIDSectionOptions struct {
	// Matches the subscriber ID as sent by the relay agent, as defined in RFC3993.
	SubscrID string `json:"subscrid,omitempty"`
	// The tag that matching clients will be assigned.
	NetworkID string `json:"networkid,omitempty"`
	// Additional options to be added for this tag aka networkid.
	DHCPOption []string `json:"dhcp_option,omitempty"`
	// Whether to send the additional options from dhcp_option list to the clients that didn't request them.
	Force uci.Bool `json:"force,omitempty"`
}

func (SubscrIDSectionOptions) IsConfigSectionOptions() {}

type TagSection struct {
	uci.StaticSectionOptions `json:",inline"`
	TagSectionOptions        `json:",inline"`
}

func (in *TagSection) DeepCopyInto(out *TagSection) {
	*out = *in
}

type TagSectionOptions struct {
	// Additional options to be added for this tag aka networkid.
	DHCPOption uci.List `json:"dhcp_option,omitempty"`
	// Whether to send the additional options from dhcp_option list to the clients that didn't request them.
	Force uci.Bool `json:"force,omitempty"`
}

func (TagSectionOptions) IsConfigSectionOptions() {}

type UserClassSection struct {
	uci.StaticSectionOptions `json:",inline"`
	UserClassSectionOptions  `json:",inline"`
}

func (in *UserClassSection) DeepCopyInto(out *UserClassSection) {
	*out = *in
}

type UserClassSectionOptions struct {
	// String sent by the client representing the user of the client. dnsmasq performs a substring match on the user
	// class string using this value.
	UserClass string `json:"userclass,omitempty"`
	// The tag that matching clients will be assigned.
	NetworkID string `json:"networkid,omitempty"`
	// Additional options to be added for this tag aka networkid.
	DHCPOption []string `json:"dhcp_option,omitempty"`
	// Whether to send the additional options from dhcp_option list to the clients that didn't request them.
	Force uci.Bool `json:"force,omitempty"`
}

func (UserClassSectionOptions) IsConfigSectionOptions() {}

type VendorClassSection struct {
	uci.StaticSectionOptions  `json:",inline"`
	VendorClassSectionOptions `json:",inline"`
}

func (in *VendorClassSection) DeepCopyInto(out *VendorClassSection) {
	*out = *in
}

type VendorClassSectionOptions struct {
	// String sent by the client representing the vendor of the client. dnsmasq performs a substring match on the
	// vendor class string using this value.
	VendorClass string `json:"vendorclass,omitempty"`
	// The tag that matching clients will be assigned.
	NetworkID string `json:"networkid,omitempty"`
	// Additional options to be added for this tag aka networkid.
	DHCPOption []string `json:"dhcp_option,omitempty"`
	// Whether to send the additional options from dhcp_option list to the clients that didn't request them.
	Force uci.Bool `json:"force,omitempty"`
}

func (VendorClassSectionOptions) IsConfigSectionOptions() {}
