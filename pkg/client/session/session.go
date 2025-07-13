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

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/rpc"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
)

type SessionInterface interface {
	Login(ctx context.Context, opts LoginOptions) (r ubus.Response, err error)
}

// implements SessionInterface
type SessionClient struct {
	call   *ubus.Call
	client rpc.UbusRPCInterface
}

func NewSessionClient(r *rpc.UbusRPCClient) (c *SessionClient) {
	c = &SessionClient{
		call: &ubus.Call{},
	}
	c.call.SetPath("session")
	c.client = r
	return c
}

func (c *SessionClient) Login(ctx context.Context, opts LoginOptions) (ubus.Response, error) {
	c.call.SetProcedure("login")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

/*
################################################################
#
# all xOptions types are in this block. they all implement the
# Signature interface.
#
################################################################
*/

// implements Signature interface
type LoginOptions struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  uint   `json:"timeout"`
}

func (LoginOptions) IsOptsType() {}

func (opts LoginOptions) GetResult(p ubus.Response) (u LoginResult, err error) {
	if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case ubus.SessionResult:
			json.Unmarshal(data, &u)
		default:
			return LoginResult{}, errors.New("not a LoginResult")
		}
	} else { // error
		return LoginResult{}, errors.New(p[0].(ubus.ExitCode).Error())
	}
	return u, nil
}

/*
################################################################
#
# all exported XResult types are in this block.
#
################################################################
*/

// result of a `session login` command
type LoginResult struct {
	session.Session `json:",inline"`
}
