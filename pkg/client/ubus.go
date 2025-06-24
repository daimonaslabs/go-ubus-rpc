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

package client

import (
	"encoding/json"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
)

type Params []any

// all implementations have an implicit method of GetResult(Response) (xResult, error)
type Signature interface {
	isOptsType()
}

// implements Signature
type Call struct {
	SessionID session.SessionID
	Path      string
	Procedure string
	Signature Signature
}

func (c *Call) asParams() Params {
	return Params{c.SessionID, c.Path, c.Procedure, c.Signature}
}

func (c *Call) setSessionID(id session.SessionID) {
	c.SessionID = id
}

func (c *Call) setPath(p string) {
	c.Path = p
}

func (c *Call) setProcedure(p string) {
	c.Procedure = p
}

func (uc *Call) setSignature(sig Signature) {
	data, err := json.Marshal(sig)
	if err != nil {
		panic(err)
	}
	uc.Signature = sig
	json.Unmarshal(data, &uc.Signature)
}
