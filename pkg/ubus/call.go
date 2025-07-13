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

package ubus

import (
	"encoding/json"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
)

type Params []any

// all implementations have an implicit method of GetResult(Response) (xResult, error)
type Signature interface {
	IsOptsType()
}

// implements Signature
type Call struct {
	sessionID session.SessionID
	path      string
	procedure string
	signature Signature
}

func (c *Call) AsParams() Params {
	return Params{c.sessionID, c.path, c.procedure, c.signature}
}

func (c *Call) SetSessionID(id session.SessionID) {
	c.sessionID = id
}

func (c *Call) SetPath(p string) {
	c.path = p
}

func (c *Call) SetProcedure(p string) {
	c.procedure = p
}

func (uc *Call) SetSignature(sig Signature) {
	data, err := json.Marshal(sig)
	if err != nil {
		panic(err)
	}
	uc.signature = sig
	json.Unmarshal(data, &uc.signature)
}
