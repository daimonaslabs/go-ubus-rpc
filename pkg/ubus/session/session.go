package session

//type SessionGetter interface {
//	Login(opts *LoginOptions) (*ubus.UbusResponse, error)
//}
//
//type sessionCall struct {
//	*client.Call
//}
//
//type Session interface {
//	Session() *UbusRPC
//}
//
//type SessionInterface interface {
//	Login(opts *.loginoptiosessionns) (*ubus.UbusResponse, error)
//}
//
//func (uc *UbusRPC) Session() {
//	uc.rpc.Call.Path = "session"
//}
//
//type sessionCall struct {
//	ubus.UbusCall
//}
//
//func (c sessionCall) Login() {
//	c.Procedure = "login"
//}

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
