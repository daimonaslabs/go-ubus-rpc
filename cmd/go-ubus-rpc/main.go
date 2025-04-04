package main

import (
	"fmt"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
)

func main() {
	sesh := client.NewSession("root", "D@!monas", "http://10.0.0.1/ubus")

	fmt.Println("sesh: ", sesh)
}
