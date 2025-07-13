# go-ubus-rpc
## Introduction
A Go library to make RPCs to OpenWrt's ubus along with a handy CLI utility, gur.

## Getting Started
### Installation
To program with it, import it like any other Go library. To use the CLI, see the releases page and download
the proper binary for your system.

### Usage
Example main.go:
```
package main

import (
	"context"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/client/uci"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/rpc"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
)

func main() {
	// create client caller
	ctx := context.Background()
	opts := rpc.UbusRPCClientOptions{Username: "root", Password: "D@!monas", URL: "http://10.0.0.1/ubus", Timeout: session.DefaultSessionTimeout}
	clientset, _ := client.NewForOpts(ctx, opts)

	// make an RPC
	uciGetOpts := uci.GetOptions{Config: "firewall"}    // declare parameters for the call
	response, _ := clientset.UCI().Get(ctx, uciGetOpts) // make the call
	result, _ := uciGetOpts.GetResult(response)         // get the typed result object from the response, in this case `result` will be a `GetResult`
}
```

To use the CLI, first do a `gur login`:
```
$ gur login -u root -p 'D@!monas' --url http://10.0.0.1/ubus
```

This will write a file to `~/.go-ubus-rpc/config.json` containing your connection info. Future gur commands
will read this file and use it to make future calls. After logging in, you can run other commands:
```
$ gur uci get -c network
{                                                                                  
  "sections": [                                                                    
    {                                                                              
      ".anonymous": false,                                                         
      ".type": "interface",                                                        
      ".name": "loopback",                                                         
      ".index": 0,                                                                 
      "device": "lo"                                                               
    },                                                                             
    {                                                                              
      ".anonymous": false,                                                         
      ".type": "globals",                                                          
      ".name": "globals",                                                          
      ".index": 1,                                                                 
      "packet_steering": "1",                                                      
      "ula_prefix": "fd29:9d35:2b57::/48"                                          
    },
    ...

$ gur uci get -c network -t device
{                                                                                  
  "sections": [                                                                    
    {                                                                              
      ".anonymous": true,                                                          
      ".type": "device",                                                           
      ".name": "cfg030f15",                                                        
      ".index": 2,                                                                 
      "name": "br-lan",                                                            
      "ports": [                                                                   
        "eth1",                                                                    
        "lan1",                                                                    
        "lan2",                                                                    
        "lan3",                                                                    
        "lan4",                                                                    
        "sfp2"                                                                     
      ],                                                                           
      "type": "bridge"                                                             
    },
    ...
```

## Description
### Overview
`Clientset` from `pkg/client/clientset` is the main client and calling object. Get one with `NewForClient`
or `NewForOpts` and use it to make calls to the remote OpenWrt instance. Each ubus command ('path' in ubus
docs) is an interface which contains all the subcommands ('procedure' in ubus docs) associated with that
top level command. The parameters for each command ('signature' in ubus docs) is also an interface because
every command has different parameters. Because of this, typed result objects are also retrieved from the
response via methods tied to the options object passed to the command.

With this design, ubus commands can be constructed in a similar fashion to how they would
be when using ubus on the command line directly.

Ubus responses are quite dynamic, the structure and content of the response changes based on the command 
executed and some commands give different responses based on the parameters passed to it. Every signature
type should implement a 'de facto' interface with a `func (opts XOptions) GetResult(p Response) (u XResult, err error)`
method which will parse the response and properly marshal it into the corresponding XResult object. This
library deals with all that bespoke JSON marshaling and unmarshaling logic so that you don't have to.

### How Commands are Constructed

All commands are built starting from a top level `Clientset` object because each command needs a ubus 
session ID and this ID is stored within this object. The command is a method on the `Clientset` object
which returns an interface containing methods which correspond to all of that command's subcommands.

### Response Handling

The `Response` object is a slice of the `ResultObject` interface, but in practice it is effectively a tuple
as ubus responses only contain two objects. The first is always the exit code of the command and the second
is the actual content of the response, the result that the user cares about. However, some commands return
with no second object in the response when successful. In that case, the exit code value (`Response[0]`) will
be zero and the error returned will be `nil`, giving the user two ways to check if the command worked.

XResult objects in `pkg/ubus` in this repo are meant to handle the raw JSON responses directly, which will
then be marshaled into an XResult type in one of `pkg/client/*` to be used by the consumer. These `/pkg/client/*` 
XResult objects aim to be more useful and easy to use for the user than the raw responses handled by `pkg/ubus`.

Each `pkg/ubus` XResult object must implement the `ResultObject` interface and may also implement an optional 
'de facto' interface with a `func matchXResult(data json.RawMessage) (ResultObject, error)` function that is
called in the `Response` object's `UnmarshalJSON` method. In order to be used in `UnmarshalJSON`, this function
must be added to the `resultObjectMatcherRegistry` via `init()`. `Response`'s `UnmarshalJSON` will then marshal
the result (`Response[1]`) into the correct xResult object, and the Signature's `GetResult` method will then
marshal the `pkg/ubus` XResult into the `pkg/client/*` XResult type.