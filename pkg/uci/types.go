package uci

type UCIConfigOptionsStatic struct {
	Anonymous bool   `json:".anonymous"`
	Type      string `json:".type"`
	Name      string `json:".name"`
	Index     int    `json:".index"`
}

var StringBoolFalse = "0"
var StringBoolTrue = "1"

type IP string
type MAC string
type StringBool string
type Time string

func ToStringBool(b bool) string {
	if b {
		return StringBoolTrue
	} else {
		return StringBoolFalse
	}
}

func True() string {
	return StringBoolTrue
}

func False() string {
	return StringBoolFalse
}

//type IP struct {
//	netip.Addr
//}
//
//// This DeepCopyInto is a manually created deepcopy function, copying the receiver, writing into out.
//// It must be non-nil.
//func (in *IP) DeepCopyInto(out *IP) {
//	*out = *in
//}
//
//// This DeepCopy is a manually created deepcopy function, copying the receiver, creating a new IP.
//func (in *IP) DeepCopy() *IP {
//	if in == nil {
//		return nil
//	}
//	out := new(IP)
//	in.DeepCopyInto(out)
//	return out
//}
//
//type MAC struct {
//	net.HardwareAddr
//}
//
//// This DeepCopyInto is a manually created deepcopy function, copying the receiver, writing into out.
//// It must be non-nil.
//func (in *MAC) DeepCopyInto(out *MAC) {
//	*out = *in
//}
//
//// This DeepCopy is a manually created deepcopy function, copying the receiver, creating a new MAC.
//func (in *MAC) DeepCopy() *MAC {
//	if in == nil {
//		return nil
//	}
//	out := new(MAC)
//	in.DeepCopyInto(out)
//	return out
//}
//
//// TODO make this safer, split into separate types which only contain required values
//// see FirewallConfig.RuleSection for all the different time and date options needed
//type Time struct {
//	time.Time
//}
//
//// This DeepCopyInto is a manually created deepcopy function, copying the receiver, writing into out.
//// It must be non-nil.
//func (in *Time) DeepCopyInto(out *Time) {
//	*out = *in
//}
//
//// This DeepCopy is a manually created deepcopy function, copying the receiver, creating a new Time.
//func (in *Time) DeepCopy() *Time {
//	if in == nil {
//		return nil
//	}
//	out := new(Time)
//	in.DeepCopyInto(out)
//	return out
//}
//
