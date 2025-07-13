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

package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus"
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

type CtxKey string

func AddToContext(ctx context.Context, u UbusRPCClient) context.Context {
	return context.WithValue(context.Background(), CtxKey("client"), u)
}

func GetFromContext(ctx context.Context) *UbusRPCClient {
	u := UbusRPCClient(ctx.Value(CtxKey("client")).(UbusRPCClient))
	return &u
}

type UbusRPCClientOptions struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  uint   `json:"timeout"`
	URL      string `json:"url"`
}
type UbusRPCInterface interface {
	Do(context.Context, *ubus.Call) (ubus.Response, error)
	Load() (string, error)
	Save()
}

type UbusRPCClient struct {
	RPCClient   *rpc.Client      `json:"-"`
	UbusSession *session.Session `json:"session"`
	URL         string           `json:"url"`
}

func (u *UbusRPCClient) Do(ctx context.Context, call *ubus.Call) (r ubus.Response, err error) {
	call.SetSessionID(u.UbusSession.SessionID)
	err = u.RPCClient.CallContext(ctx, &r, "call", call.AsParams()...)
	if len(r) == 0 {
		err = errors.New("empty response")
		return nil, err
	}
	if r[0].(ubus.ExitCode) != 0 {
		err = errors.New(r[0].(ubus.ExitCode).Error())
	}
	return r, err
}

func (u *UbusRPCClient) Load() (string, error) {
	configFileBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return configFilePath, err
	} else {
		err = json.Unmarshal(configFileBytes, &u)
		if err != nil {
			log.Fatalln(err)
		}

		u.RPCClient, err = newRPCClient(context.Background(), u.URL)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return configFilePath, nil
}

func (u *UbusRPCClient) Save() {
	configFileBytes, err := json.MarshalIndent(u, "", "  ")

	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(configFilePath, configFileBytes, 0600)
	if err != nil {
		log.Fatalln(err)
	}
}

func NewUbusRPCClient(ctx context.Context, opts UbusRPCClientOptions) (*UbusRPCClient, error) {
	// initialize RPC client
	rpcClient, err := newRPCClient(ctx, opts.URL)
	if err != nil {
		log.Fatalln(err)
	}

	c := UbusRPCClient{
		RPCClient:   rpcClient,
		UbusSession: &session.Session{},
		URL:         opts.URL,
	}
	loginOpts := loginOptions{
		Username: opts.Username,
		Password: opts.Password,
		Timeout:  opts.Timeout,
	}
	login := &ubus.Call{}
	login.SetSessionID(session.LoginSessionID)
	login.SetPath("session")
	login.SetProcedure("login")
	login.SetSignature(loginOpts)

	// initialize ubus client
	response := ubus.Response{}
	if err != nil {
		log.Fatalln(err)
	}

	err = c.RPCClient.CallContext(ctx, &response, "call", login.AsParams()...)
	if err != nil {
		log.Fatalln(err)
	}

	session := response[1].(ubus.SessionResult)
	c.UbusSession = &session.Session
	return &c, err
}

func newRPCClient(ctx context.Context, url string) (*rpc.Client, error) {
	tokenHeader := rpc.WithHeader("Content-Type", "application/json")
	httpClient := rpc.WithHTTPClient(&http.Client{
		Timeout: 10 * time.Second,
	})
	return rpc.DialOptions(ctx, url, httpClient, tokenHeader)
}

// implements Signature interface
type loginOptions struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  uint   `json:"timeout"`
}

func (loginOptions) IsOptsType() {}
