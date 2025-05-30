package uci

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
)

func NewUCICommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "uci",
		Short: "Run UCI commands.",
		Long:  "Run UCI commands to update router configs.",
		PersistentPreRun: func(c *cobra.Command, args []string) {
			rpc := client.UbusRPC{}
			configFile, err := rpc.Load()
			if err != nil {
				log.Fatalln(configFile, "not found, you should run `gur login`!")
			}
			ctx := client.AddToContext(c.Context(), rpc)
			c.SetContext(ctx)
		},
	}

	c.AddCommand(
		NewGetCommand(),
	)

	return c
}

func NewGetCommand() *cobra.Command {
	o := GetOptions{}
	structType := reflect.TypeOf(o)
	numOptions := structType.NumField()
	c := &cobra.Command{
		Use:   "get",
		Short: "Get a config value",
		Args:  cobra.MaximumNArgs(numOptions),
		RunE: func(c *cobra.Command, args []string) error {
			return o.Run(c)
		},
	}
	o.BindFlags(c)

	return c
}

type GetOptions struct {
	Config  string
	Section string
	Type    string
	Option  string
}

func (o *GetOptions) BindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.Config, "config", "c", "", "Which config to query.")
	c.Flags().StringVarP(&o.Section, "section", "s", "", "The section of the config.")
	c.Flags().StringVarP(&o.Type, "type", "t", "", "The type of the config section.")
	c.Flags().StringVarP(&o.Option, "option", "o", "", "A single option within a config section.")
	c.MarkFlagRequired("config")
}

func (o *GetOptions) Run(c *cobra.Command) error {
	uciGetOpts := client.UCIGetOptions{
		Config:  uci.ConfigName(o.Config),
		Section: o.Section,
		Type:    o.Type,
		Option:  o.Option,
	}
	ctx := c.Context()
	rpc := client.GetFromContext(c.Context())
	response, err := rpc.UCI().Get(ctx, uciGetOpts)
	if err != nil {
		return err
	}
	result, err := uciGetOpts.GetResult(response)
	if err != nil {
		return err
	}
	output, err := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))
	return err
}
