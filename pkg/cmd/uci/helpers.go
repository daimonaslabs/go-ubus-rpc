/*
Copyright 2025 Daimonas Labs.

Licensed under the GNU General Public License, Version 3 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.gnu.org/licenses/gpl-3.0.en.html

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
