package uci

type UCIConfigOptionsStatic struct {
	Anonymous bool   `json:".anonymous"`
	Type      string `json:".type"`
	Name      string `json:".name"`
	Index     int    `json:".index"`
}

var StringBoolFalse = "0"
var StringBoolTrue = "1"

// TODO make these safer, check for the right format within the strings. split into
// separate, more specific types as needed. see FirewallConfig.RuleSection for all the
// different time and date options needed.
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
