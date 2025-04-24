package ubus

import (
	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
)

type UbusInterface interface {
	UbusClient() client.Caller
	SessionInterface() session.SessionInterface
	// UCICallGetter
}

type UbusClient struct {
	client client.Caller
}

//type UCICallGetter interface {
//	UCI() *UbusCall
//}
