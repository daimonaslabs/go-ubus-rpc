package encoding

import (
	"encoding/json"
)

// TODO unnecessary! replace all `json` tags with their `ubus` tag values, anyone else can convert to
// different names
type UbusMarshalUnmarshaler interface {
	json.Marshaler
	json.Unmarshaler
}
