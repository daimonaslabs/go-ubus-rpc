package uci

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
)

func checkConfig(c string) error {
	if !slices.Contains(uci.Configs, c) {
		return fmt.Errorf("invalid config option, must be one of: %s", uci.Configs)
	} else {
		return nil
	}
}

func unmarshalCLIValues[S uci.ConfigSectionOptions](o *SetOptions) (u client.UCISetOptions) {
	var s S
	err := json.Unmarshal([]byte(o.Values), &s)
	if err != nil {
		log.Fatalln(err)
	}
	u = client.UCISetOptions{
		Config:  o.Config,
		Section: o.Section,
		Values:  s,
	}
	return u
}
