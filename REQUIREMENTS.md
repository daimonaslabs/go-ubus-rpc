# Requirements

## Libraries
#### [ethereum-go/rpc](https://pkg.go.dev/github.com/ethereum/go-ethereum/rpc)

## Implementation
- Separate request and response types (request is struct, response is interface)
    - net/http will be involved somehow
- Custom marshaller for \`ubus\` tags ([Config.TagKey](https://pkg.go.dev/github.com/json-iterator/go#Config))
    - Marshal bools to 1 or 0
- ubus `Path` = Go type, `Procedure` = Methods

Example:
```
type UbusResponse interface {
    Write()
}
type UbusRequest interface {
    String()
}
type UCIRequest struct {
   UbusSessionId 
}
func (u UCIRequest) String() {
    return ...
}

type UCIResponse struct {
    ...
}

type UCI struct {
    rpc.Client
}
func (u UCI) Get(UCIRequest) UCIResponse {
    ...
}
```