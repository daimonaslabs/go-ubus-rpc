package v1alpha1

type UCIConfigOptionsStatic struct {
	Anonymous bool   `json:"anonymous" ubus:".anonymous"`
	Type      string `json:"type" ubus:".type"`
	Name      string `json:"name" ubus:".name"`
	Index     int    `json:"index" ubus:".index"`
}
