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
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
	"github.com/ethereum/go-ethereum/rpc"
)

var configDirPath, configFilePath string

func init() {
	configDirPath = os.Getenv("HOME") + "/.go-ubus-rpc"
	configFilePath = filepath.Join(configDirPath, "config.json")
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		err := os.MkdirAll(configDirPath, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

type clientset struct {
	RPCClient   *rpc.Client      `json:"-"`
	UbusSession *session.Session `json:"session"`
	URL         string           `json:"url"`
}

// the primary client and caller object
type UbusRPC struct {
	Call
	clientset
	sessionRPC
	uciRPC
}

func (u *UbusRPC) Session() SessionInterface {
	return newSessionRPC(u)
}

func (u *UbusRPC) UCI() UCIInterface {
	return newUCIRPC(u)
}

type ClientOptions struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  uint   `json:"timeout"`
	URL      string `json:"url"`
}

func newRPCClient(ctx context.Context, url string) (*rpc.Client, error) {
	tokenHeader := rpc.WithHeader("Content-Type", "application/json")
	httpClient := rpc.WithHTTPClient(&http.Client{
		Timeout: 10 * time.Second,
	})
	return rpc.DialOptions(ctx, url, httpClient, tokenHeader)
}

func NewUbusRPC(ctx context.Context, opts *ClientOptions) (*UbusRPC, error) {
	// initialize RPC client
	rpcClient, err := newRPCClient(ctx, opts.URL)
	if err != nil {
		log.Fatalln(err)
	}

	c := clientset{
		RPCClient:   rpcClient,
		UbusSession: &session.Session{},
		URL:         opts.URL,
	}
	loginOpts := SessionLoginOptions{
		Username: opts.Username,
		Password: opts.Password,
		Timeout:  opts.Timeout,
	}
	login := &Call{}
	login.setSessionID(session.LoginSessionID)
	login.setPath("session")
	login.setProcedure("login")
	login.setSignature(loginOpts)

	// initialize ubus client
	response := Response{}
	if err != nil {
		log.Fatalln(err)
	}

	err = c.RPCClient.CallContext(ctx, &response, "call", login.asParams()...)
	if err != nil {
		log.Fatalln(err)
	}

	session := response[1].(sessionResult)
	c.UbusSession = &session.Session
	return &UbusRPC{
		Call: Call{
			SessionID: c.UbusSession.SessionID,
		},
		clientset: c,
	}, err
}

func (u *UbusRPC) do(ctx context.Context) (r Response, err error) {
	err = u.clientset.RPCClient.CallContext(ctx, &r, "call", u.Call.asParams()...)
	if len(r) == 0 {
		err = errors.New("empty response")
		return nil, err
	}
	if r[0].(ExitCode) != 0 {
		err = errors.New(r[0].(ExitCode).Error())
	}
	return r, err
}

func (u *UbusRPC) Save() {
	configFileBytes, err := json.MarshalIndent(u.clientset, "", "  ")

	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(configFilePath, configFileBytes, 0600)
	if err != nil {
		log.Fatalln(err)
	}
}

func (u *UbusRPC) Load() (string, error) {
	configFileBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return configFilePath, err
	} else {
		err = json.Unmarshal(configFileBytes, &u.clientset)
		if err != nil {
			log.Fatalln(err)
		}
		u.Call = Call{
			SessionID: u.UbusSession.SessionID,
		}
		u.clientset.RPCClient, err = newRPCClient(context.Background(), u.URL)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return configFilePath, nil
}

type CtxKey string

func AddToContext(ctx context.Context, u UbusRPC) context.Context {
	return context.WithValue(context.Background(), CtxKey("client"), u)
}

func GetFromContext(ctx context.Context) *UbusRPC {
	u := UbusRPC(ctx.Value(CtxKey("client")).(UbusRPC))
	return &u
}
