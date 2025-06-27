package system

import (
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
)

const (
	// the name of this config
	Config = "system"

	// these are static values for the uci.StaticSectionOptions.Type field
	System     = "system"
	Timeserver = "timeserver"
)

var (
	Sections []string
)

func init() {
	Sections = []string{System, Timeserver}
}

type SystemSection struct {
	uci.StaticSectionOptions `json:",inline"`
	SystemSectionOptions     `json:",inline"`
}

func (in *SystemSection) DeepCopyInto(out *SystemSection) {
	*out = *in
}

type SystemSectionOptions struct {
	// Kernel message buffer size.
	Buffersize *uci.Int `json:"buffersize,omitempty"`
	// Max log level for kernel messages to console (1–8).
	Conloglevel *uci.Int `json:"conloglevel,omitempty"`
	// Min level for cron messages (0–9+).
	Cronloglevel *uci.Int `json:"cronloglevel,omitempty"`
	// Short, single-line human-readable system description.
	Description *string `json:"description,omitempty"`
	// System hostname (avoid dots).
	Hostname *string `json:"hostname,omitempty"`
	// Same as conloglevel, but takes precedence.
	Klogconloglevel *uci.Int `json:"klogconloglevel,omitempty"`
	// Size of log buffer used by `logread`.
	LogBufferSize *uci.Int `json:"log_buffer_size,omitempty"`
	// Path to log file (optional).
	LogFile *string `json:"log_file,omitempty"`
	// Hostname sent to remote syslog.
	LogHostname *string `json:"log_hostname,omitempty"`
	// IP address of remote syslog server.
	LogIP *string `json:"log_ip,omitempty"`
	// Port for remote syslog server (default: 514).
	LogPort *uci.Int `json:"log_port,omitempty"`
	// Prefix for network log messages.
	LogPrefix *string `json:"log_prefix,omitempty"`
	// Protocol: "tcp" or "udp" (default: udp).
	LogProto *string `json:"log_proto,omitempty"`
	// Enable remote logging (default: true).
	LogRemote *uci.Bool `json:"log_remote,omitempty"`
	// Log buffer size in KiB (default: 64).
	LogSize *uci.Int `json:"log_size,omitempty"`
	// Use `\0` instead of `\n` with TCP.
	LogTrailerNull *uci.Bool `json:"log_trailer_null,omitempty"`
	// "circular" or "file".
	LogType *string `json:"log_type,omitempty"`
	// Freeform multiline notes (e.g. location, inventory).
	Notes *string `json:"notes,omitempty"`
	// Require login on console access (default: false).
	Ttylogin *uci.Bool `json:"ttylogin,omitempty"`
	// Path to urandom seed.
	UrandomSeed *string `json:"urandom_seed,omitempty"`
	// POSIX.1 timezone string (e.g. UTC).
	Timezone *string `json:"timezone,omitempty"`
	// IANA/Olson timezone string (e.g. Europe/London).
	Zonename *string `json:"zonename,omitempty"`
	// Compression algorithm for ZRAM (e.g. lzo, lz4, zstd).
	ZramCompAlgo *string `json:"zram_comp_algo,omitempty"`
	// ZRAM size in MB.
	ZramSizeMB *uci.Int `json:"zram_size_mb,omitempty"`
}

func (SystemSectionOptions) IsConfigSectionOptions() {}

type TimeserverSection struct {
	uci.StaticSectionOptions `json:",inline"`
	TimeserverSectionOptions `json:",inline"`
}

func (in *TimeserverSection) DeepCopyInto(out *TimeserverSection) {
	*out = *in
}

type TimeserverSectionOptions struct {
	Enabled      *uci.Bool `json:"enabled,omitempty"`
	EnableServer *uci.Bool `json:"enable_server,omitempty"`
	// List of ntp servers to query
	Server *uci.List `json:"server,omitempty"`
}

func (TimeserverSectionOptions) IsConfigSectionOptions() {}
