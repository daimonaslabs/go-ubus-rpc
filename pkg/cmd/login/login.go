package login

import (
	"context"
	"log"
	"reflect"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
	"github.com/spf13/cobra"
)

func NewLoginCommand() *cobra.Command {
	o := LoginOptions{}
	structType := reflect.TypeOf(o)
	numOptions := structType.NumField()
	c := &cobra.Command{
		Use:   "login",
		Short: "Log into ubus.",
		Long:  "Create a ubus session to make subsequent RPCs with.",
		Args:  cobra.MaximumNArgs(numOptions),
		RunE: func(c *cobra.Command, args []string) error {
			return o.Run(c)
		},
	}
	o.BindFlags(c)

	return c
}

type LoginOptions struct {
	client.ClientOptions
}

func (o *LoginOptions) BindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.URL, "url", "", "", "URL.")
	c.Flags().StringVarP(&o.Username, "username", "u", "", "Username.")
	c.Flags().StringVarP(&o.Password, "password", "p", "", "Password.")
	c.Flags().UintVarP(&o.Timeout, "timeout", "t", 0, "How long the session will last. Default (0) has no expiry.")
	c.MarkFlagRequired("url")
	c.MarkFlagRequired("username")
	c.MarkFlagRequired("password")
}

func (o *LoginOptions) Run(c *cobra.Command) error {
	ctx := context.Background()
	rpc, err := client.NewUbusRPC(ctx, &o.ClientOptions)
	if err != nil {
		log.Fatalln("error creating ubus client")
	}

	rpc.Save()
	return nil
}
