package main

import (
	"os"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/cmd"
)

func main() {
	ubus := cmd.NewCommand()
	err := ubus.Execute()
	if err != nil {
		ubus.Usage()
		os.Exit(1)
	}
}
