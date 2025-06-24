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

package session

const (
	LoginSessionID         SessionID = "00000000000000000000000000000000"
	DefaultSessionTimeout  uint      = 300
	NoExpirySessionTimeout uint      = 0
)

// maybe use this to do validation on the SessionID
//type sessionID [32]byte

type SessionID string

type Session struct {
	SessionID SessionID `json:"ubus_rpc_session"`
	Timeout   int       `json:"timeout"`
	Expires   int       `json:"expires"`
	ACLs      ACL       `json:"acls"`
	Data      Data      `json:"data"`
}

type ACL struct {
	AccessGroup map[string][]string `json:"access-group"`
	CGIIO       map[string][]string `json:"cgi-io,omitempty"`
	File        map[string][]string `json:"file,omitempty"`
	Ubus        map[string][]string `json:"ubus"`
	UCI         map[string][]string `json:"uci,omitempty"`
}

type Data struct {
	Username string `json:"username"`
}
