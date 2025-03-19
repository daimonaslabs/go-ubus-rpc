package v1alpha1

type UCIConfigOptionsStatic struct {
	Anonymous bool   `json:"anonymous" ubus:".anonymous"`
	Type      string `json:"type" ubus:".type"`
	Name      string `json:"name" ubus:".name"`
	Index     int    `json:"index" ubus:".index"`
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
