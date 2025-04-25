package main

import (
	"fmt"
	"log"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
)

func main() {
	opts := client.ClientOptions{Username: "root", Password: "D@!monas", URL: "http://10.0.0.1/ubus", Timeout: client.DefaultSessionTimeout}
	rpc, err := client.NewUbusRPC(&opts)
	login := client.LoginOptions{
		Username: opts.Username,
		Password: opts.Password,
		Timeout:  opts.Timeout,
	}
	response := rpc.Session().Login(&login)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(response)

}
