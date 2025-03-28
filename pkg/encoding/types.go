package encoding

import (
	"net"
	"net/netip"
	"time"
)

/*
This package contains types and functions which make it easy to construct ubus calls by providing type
safety to the unique ubus parameters as well as providing a new ubus encoder so that standard Go structs
can be marshaled into these types.
*/

type UCIConfigOptionsStatic struct {
	Anonymous bool   `json:"dotAnonymous" ubus:".anonymous"`
	Type      string `json:"dotType" ubus:".type"`
	Name      string `json:"dotName" ubus:".name"`
	Index     int    `json:"dotIndex" ubus:".index"`
}

var UbusBoolFalse = "0"
var UbusBoolTrue = "1"

type UbusBool string

func ConvertBool(b bool) string {
	if b {
		return UbusBoolTrue
	} else {
		return UbusBoolFalse
	}
}

func True() string {
	return UbusBoolTrue
}

func False() string {
	return UbusBoolFalse
}

type IP struct {
	netip.Addr
}

// This DeepCopyInto is a manually created deepcopy function, copying the receiver, writing into out.
// It must be non-nil.
func (in *IP) DeepCopyInto(out *IP) {
	*out = *in
}

// This DeepCopy is a manually created deepcopy function, copying the receiver, creating a new IP.
func (in *IP) DeepCopy() *IP {
	if in == nil {
		return nil
	}
	out := new(IP)
	in.DeepCopyInto(out)
	return out
}

type MAC struct {
	net.HardwareAddr
}

// This DeepCopyInto is a manually created deepcopy function, copying the receiver, writing into out.
// It must be non-nil.
func (in *MAC) DeepCopyInto(out *MAC) {
	*out = *in
}

// This DeepCopy is a manually created deepcopy function, copying the receiver, creating a new MAC.
func (in *MAC) DeepCopy() *MAC {
	if in == nil {
		return nil
	}
	out := new(MAC)
	in.DeepCopyInto(out)
	return out
}

// TODO make this safer, splitu into separate types which only contain required values
// see FirewallConfig.RuleSection for all the different time and date options needed
type Time struct {
	time.Time
}

// This DeepCopyInto is a manually created deepcopy function, copying the receiver, writing into out.
// It must be non-nil.
func (in *Time) DeepCopyInto(out *Time) {
	*out = *in
}

// This DeepCopy is a manually created deepcopy function, copying the receiver, creating a new Time.
func (in *Time) DeepCopy() *Time {
	if in == nil {
		return nil
	}
	out := new(Time)
	in.DeepCopyInto(out)
	return out
}
