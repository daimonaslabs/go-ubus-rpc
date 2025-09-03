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

package file

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/rpc"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus"
)

type FileInterface interface {
	Read(ctx context.Context, opts ReadOptions) (r ubus.Response, err error)
	Write(ctx context.Context, opts WriteOptions) (r ubus.Response, err error)
	List(ctx context.Context, opts ListOptions) (r ubus.Response, err error)
	Stat(ctx context.Context, opts StatOptions) (r ubus.Response, err error)
	MD5(ctx context.Context, opts MD5Options) (r ubus.Response, err error)
	Remove(ctx context.Context, opts RemoveOptions) (r ubus.Response, err error)
	Exec(ctx context.Context, opts ExecOptions) (r ubus.Response, err error)
}

// implements FileInterface
type FileClient struct {
	call   *ubus.Call
	client rpc.UbusRPCInterface
}

func NewFileClient(r *rpc.UbusRPCClient) (c *FileClient) {
	c = &FileClient{
		call: &ubus.Call{},
	}
	c.call.SetPath("file")
	c.client = r
	return c
}

func (c *FileClient) Read(ctx context.Context, opts ReadOptions) (ubus.Response, error) {
	c.call.SetProcedure("read")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

func (c *FileClient) Write(ctx context.Context, opts WriteOptions) (ubus.Response, error) {
	c.call.SetProcedure("write")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

func (c *FileClient) List(ctx context.Context, opts ListOptions) (ubus.Response, error) {
	c.call.SetProcedure("list")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

func (c *FileClient) Stat(ctx context.Context, opts StatOptions) (ubus.Response, error) {
	c.call.SetProcedure("stat")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

func (c *FileClient) MD5(ctx context.Context, opts MD5Options) (ubus.Response, error) {
	c.call.SetProcedure("md5")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

func (c *FileClient) Exec(ctx context.Context, opts ExecOptions) (ubus.Response, error) {
	c.call.SetProcedure("exec")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

func (c *FileClient) Remove(ctx context.Context, opts RemoveOptions) (ubus.Response, error) {
	c.call.SetProcedure("remove")
	c.call.SetSignature(opts)

	return c.client.Do(ctx, c.call)
}

/*
################################################################
#
# all XOptions types are in this block. they all implement the
# Signature interface.
#
################################################################
*/

// implements Signature interface
type ReadOptions struct {
	Path   string `json:"path,omitempty"`
	Base64 bool   `json:"base64,omitempty"`
}

func (ReadOptions) IsOptsType() {}

type ReadResult struct {
	Data string `json:"data"`
}

func (ReadResult) IsResultObject() {}

func (opts ReadOptions) GetResult(p ubus.Response) (u ReadResult, err error) {
	if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case ubus.DataResult:
			json.Unmarshal(data, &u)
		default:
			return ReadResult{}, errors.New("not a ReadResult")
		}
	} else { // error
		return ReadResult{}, errors.New(p[0].(ubus.ExitCode).Error())
	}
	return u, nil
}

// implements Signature interface
type WriteOptions struct {
	Path   string `json:"path"`
	Data   string `json:"data"`
	Append bool   `json:"append"`
	Mode   int    `json:"mode,omitempty"`
	Base64 bool   `json:"base64,omitempty"`
}

func (WriteOptions) IsOptsType() {}

// implements Signature interface
type ListOptions struct {
	Path string `json:"path"`
}

func (ListOptions) IsOptsType() {}

type ListResult struct {
	Entries []ubus.FileResult `json:"entries"`
}

func (ListResult) IsResultObject() {}

func (opts ListOptions) GetResult(p ubus.Response) (u ListResult, err error) {
	if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case ubus.EntriesResult:
			json.Unmarshal(data, &u)
		default:
			return ListResult{}, errors.New("not a ListResult")
		}
	} else { // error
		return ListResult{}, errors.New(p[0].(ubus.ExitCode).Error())
	}
	return u, nil
}

// implements Signature interface
type StatOptions struct {
	Path string `json:"path"`
}

func (StatOptions) IsOptsType() {}

func (opts StatOptions) GetResult(p ubus.Response) (u ubus.FileResult, err error) {
	if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case ubus.FileResult:
			json.Unmarshal(data, &u)
		default:
			return ubus.FileResult{}, errors.New("not a FileResult")
		}
	} else { // error
		return ubus.FileResult{}, errors.New(p[0].(ubus.ExitCode).Error())
	}
	return u, nil
}

// implements Signature interface
type MD5Options struct {
	Path string `json:"path"`
}

func (MD5Options) IsOptsType() {}

type MD5Result struct {
	MD5 string `json:"md5"`
}

func (MD5Result) IsResultObject() {}

func (opts MD5Options) GetResult(p ubus.Response) (u MD5Result, err error) {
	if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case ubus.MD5Result:
			json.Unmarshal(data, &u)
		default:
			return MD5Result{}, errors.New("not a MD5Result")
		}
	} else { // error
		return MD5Result{}, errors.New(p[0].(ubus.ExitCode).Error())
	}
	return u, nil
}

// implements Signature interface
type RemoveOptions struct {
	Path string `json:"path"`
}

func (RemoveOptions) IsOptsType() {}

// implements Signature interface
type ExecOptions struct {
	Command string   `json:"command"`
	Params  []string `json:"params"`
}

func (ExecOptions) IsOptsType() {}

func (opts ExecOptions) GetResult(p ubus.Response) (u MD5Result, err error) {
	if len(p) > 1 {
		data, _ := json.Marshal(p[1])
		switch p[1].(type) {
		case ubus.MD5Result:
			json.Unmarshal(data, &u)
		default:
			return MD5Result{}, errors.New("not a MD5Result")
		}
	} else { // error
		return MD5Result{}, errors.New(p[0].(ubus.ExitCode).Error())
	}
	return u, nil
}
