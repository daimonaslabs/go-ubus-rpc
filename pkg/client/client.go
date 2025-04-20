package client

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/session"
	"github.com/ethereum/go-ethereum/rpc"
)

type SessionManager struct {
	rpcClient   *rpc.Client
	ubusSession *session.Session
	options     SessionManagerOptions
}

type SessionManagerOptions struct {
	KeepAlive int
}

func NewSessionManager(url string, opts *SessionManagerOptions) *SessionManager {
	ctx := context.Background()
	tokenHeader := rpc.WithHeader("Content-Type", "application/json")
	httpClient := rpc.WithHTTPClient(&http.Client{
		Timeout: 10 * time.Second,
	})

	rpcClient, err := rpc.DialOptions(ctx, url, httpClient, tokenHeader)

	if err != nil {
		log.Fatalln(err)
	}

	return &SessionManager{
		rpcClient,
		&session.Session{},
		*opts,
	}
}

// equivalent of http.Server.ListenAndServe()
func (sm *SessionManager) ManageSession(*session.Session) {
	// TODO renew session every KeepAlive period
}

func (sm *SessionManager) GetSessionID() ubus.SessionID {
	return sm.ubusSession.SessionID
}
