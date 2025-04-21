package main

import (
	"context"
	"fmt"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/client/ubus"
)

func main() {
	ctx := context.Background()
	sesh := client.NewClient(ctx, &client.ClientOptions{Username: "root", Password: "D@!monas", URL: "http://10.0.0.1/ubus", Timeout: ubus.DefaultSessionTimeout})

	fmt.Println("sesh: ", sesh)
}
