package session

import (
	"github.com/daimonaslabs/go-ubus-rpc/pkg/client/ubus"
)

type SessionCaller interface {
	Login(opts *LoginOptions)
}

func Login(us *ubus.Session, opts *LoginOptions) {

}

type LoginOptions struct {
	Username string
	Password string
	Timeout  int
}

/*
type SessionLoginHandler struct {
	*ubus.UbusCall
	ubus.UbusResponseWriter
}

type SessionLoginResponse struct {
	Session
}

func (slr SessionLoginResponse) Write(resp []byte) error {
	resp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(resp, slr)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}
*/
