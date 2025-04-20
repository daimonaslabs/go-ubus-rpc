package ubus

type UbusHandler interface {
	UbusRPC(*UbusCall, UbusResponseWriter) error
}

type UbusHandlerFunc func(*UbusCall, UbusResponseWriter) error

func (uhf UbusHandlerFunc) UbusRPC(uc *UbusCall, urw UbusResponseWriter) error {
	return uhf(uc, urw)
}

type SessionID string
type Path string
type Procedure string
type Signature map[string]any

type UbusCall struct {
	SessionID SessionID
	Path      Path
	Procedure Procedure
	Signature Signature
}

type UbusResponseWriter interface {
	Write([]byte) error
}
