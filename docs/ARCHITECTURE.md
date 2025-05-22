# Architecture

## Overview
`UbusRPC` is the main client and calling object. Get one with `NewUbusRPC` and use it to make calls to
the remote OpenWRT instance. Each ubus command ('path' in ubus docs) is an interface which contains
all the subcommands ('procedure' in ubus docs) associated with that top level command. The parameters
for each command ('signature' in ubus docs) is also an interface because every command has different
parameters. Because of this, typed result objects are also retrieved from the response via methods tied to
the options object passed to the command.

With this design, ubus commands can be constructed in a similar fashion to how they would
be when using ubus on the command line directly.

For example:
```
package main

import (
	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
)

func main() {
// create client caller
clientOpts := client.ClientOptions{Username: "root", Password: "D@!monas", URL: "http://10.0.0.1/ubus", Timeout: session.DefaultSessionTimeout}
rpc, _ := client.NewUbusRPC(ctx, &clientOpts)

// make an RPC
uciGetOpts := client.UCIGetOptions{Config: "firewall"} // declare parameters for the call
response, _ := rpc.UCI().Get(uciGetOpts)               // make the call
result, _ := uciGetOpts.GetResult(response)            // get the typed result object from the response, in this case `result` will be a `UCIGetResult`
}
```

Ubus responses are quite dynamic, the structure and content of the response changes based on the command 
executed and some commands give different responses based on the parameters passed to it. Every signature
type should implement a 'de facto' interface with a `func (opts xOptions) GetResult(p Response) (u xResult, err error)`
method which will parse the response and properly marshal it into the corresponding xResult object. This
library deals with all that bespoke JSON marshaling and unmarshaling logic so that you don't have to.

## How Commands are Constructed

All commands are built starting from a top level `UbusRPC` object because each command needs a ubus 
session ID and this ID is stored within this object. The command is a method on the `UbusRPC` object
which returns an interface containing methods which correspond to all of that command's subcommands.
This interface is implemented by an unexported xRPC type which embeds `*UbusRPC`. 

For example:
```
func (u *UbusRPC) UCI() UCIInterface {
	return newUCIRPC(u)
}

func newUCIRPC(u *UbusRPC) *uciRPC {
	u.Call.setPath("uci")
	return &uciRPC{u}
}

type UCIInterface interface {
	Get(ctx context.Context, opts UCIGetOptions) (r Response, err error)
    ...
}

// implements UCIInterface
type uciRPC struct {
	*UbusRPC
}

func (c *uciRPC) Get(ctx context.Context, opts UCIGetOptions) (Response, error) {
	c.setProcedure("get")
	c.setSignature(opts)

	return c.do(ctx)
}
```

## Response Handling

The `Response` object is a slice of the `ResultObject` interface, but in practice it is effectively a tuple
as ubus responses only contain two objects. The first is always the exit code of the command and the second
is the actual content of the response, the result that the user cares about. However, some commands return
with no second object in the response when successful. In that case, the exit code value (`Response[0]`) will
be zero and the error returned will be `nil`, giving the user two ways to check if the command worked.

Unexported xResult objects in this repo are meant to handle the raw JSON responses directly, which will
then be marshaled into an exported XResult type to be used by the consumer. These exported XResult objects
aim to be more useful and easy to use for the user than the raw responses handled by the unexported xResult
objects.

Each unexported xResult object must implement the `ResultObject` interface and may also implement an optional 
'de facto' interface with a `func matchxResult(data json.RawMessage) (ResultObject, error)` function that is
called in the `Response` object's `UnmarshalJSON` method. In order to be used in `UnmarshalJSON`, this function
must be added to the `resultObjectMatcherRegistry` via `init()`. `Response`'s `UnmarshalJSON` will then marshal
the result (`Response[1]`) into the correct xResult object, and the Signature's `GetResult` method will then
marshal the unexported xResult into the exported XResult type.

For example:
```
// implements Signature interface
type UCIGetOptions struct {
	Config  string `json:"config,omitempty"`
	Section string `json:"section,omitempty"`
	Type    string `json:"type,omitempty"`
	Option  string `json:"option,omitempty"`
}

func (UCIGetOptions) isOptsType() {}

func (opts UCIGetOptions) GetResult(p Response) (u UCIGetResult, err error) {
	if len(p) > 1 {
		switch obj := p[1].(type) {
		case valueResult:
			u.Option = map[string]string{opts.Option: obj.Value}
		case valuesResult:
			...
		default:
			return u, errors.New("not a GetResult")
		}
	} else { // error
		return u, errors.New(p[0].(ExitCode).Error())
	}
	return u, err
}

// result of a `uci get` command
type UCIGetResult struct {
	SectionArray []uci.UCIConfigSection `json:"sectionArray,omitempty"`
	Option       map[string]string      `json:"option,omitempty"`
}

type valuesResult struct {
	Values map[string]uci.UCIConfigSection `json:"values"`
}

func (valuesResult) isResultObject() {}

func (v valuesResult) MarshalJSON() ([]byte, error) {
	...
	return json.Marshal(out)
}

func (v *valuesResult) UnmarshalJSON(data []byte) (err error) {
	...
	return nil
}
```