package ubus

import (
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
)

type Ubus interface {
	Session() session.SessionCaller
}

type UbusCaller interface {
	Session() session.SessionCaller
}
