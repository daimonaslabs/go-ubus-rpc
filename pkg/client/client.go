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
	"context"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client/file"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/client/session"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/client/uci"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/rpc"
)

type Interface interface {
	UbusRPCClient() rpc.UbusRPCInterface
	Session() session.SessionInterface
	UCI() uci.UCIInterface
}

// the primary client and caller object
type Clientset struct {
	rpcClient     rpc.UbusRPCClient
	sessionClient *session.SessionClient
	uciClient     *uci.UCIClient
}

func (u *Clientset) UbusRPCClient() rpc.UbusRPCInterface {
	return &u.rpcClient
}

func (u *Clientset) Session() session.SessionInterface {
	return session.NewSessionClient((&u.rpcClient))
}

func (u *Clientset) UCI() uci.UCIInterface {
	return uci.NewUCIClient(&u.rpcClient)
}

func (u *Clientset) File() file.FileInterface {
	return file.NewFileClient(&u.rpcClient)
}

// create a new clientset for the given RPC client
func NewForClient(ctx context.Context, client *rpc.UbusRPCClient) (c *Clientset, err error) {
	return &Clientset{
		rpcClient: *client,
	}, nil
}

// create a new clientset for the given RPC client options
func NewForOpts(ctx context.Context, opts rpc.UbusRPCClientOptions) (c *Clientset, err error) {
	urc, err := rpc.NewUbusRPCClient(ctx, opts)
	return &Clientset{
		rpcClient: *urc,
	}, err
}
