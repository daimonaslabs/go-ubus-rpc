package cmd

import (
	"github.com/spf13/cobra"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/cmd/login"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/cmd/uci"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "gur",
		Short: "Make RPCs to OpenWRT's ubus.",
		Long:  "Make RPCs to OpenWRT's ubus.",
	}

	c.AddCommand(
		login.NewLoginCommand(),
		uci.NewUCICommand(),
	)

	return c
}
