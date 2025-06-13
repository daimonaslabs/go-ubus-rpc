package uci

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/daimonaslabs/go-ubus-rpc/pkg/client"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/dhcp"
	"github.com/daimonaslabs/go-ubus-rpc/pkg/ubus/uci/firewall"
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
		NewAddCommand(),
		NewApplyCommand(),
		NewChangesCommand(),
		NewConfigsCommand(),
		NewDeleteCommand(),
		NewGetCommand(),
		NewRevertCommand(),
		NewSetCommand(),
	)

	return c
}

func NewAddCommand() *cobra.Command {
	o := AddOptions{}
	structType := reflect.TypeOf(o)
	numOptions := structType.NumField()
	c := &cobra.Command{
		Use:   "add",
		Short: "Add a new config section.",
		Args:  cobra.MaximumNArgs(numOptions),
		RunE: func(c *cobra.Command, args []string) error {
			return o.Run(c)
		},
	}
	o.BindFlags(c)

	return c
}

type AddOptions struct {
	Config string
	Type   string
}

func (o *AddOptions) BindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.Config, "config", "c", "", "Which config to query.")
	c.Flags().StringVarP(&o.Type, "type", "t", "", "The type of the config section.")
	c.MarkFlagRequired("config")
	c.MarkFlagRequired("type")
}

func (o *AddOptions) Run(c *cobra.Command) (err error) {
	if err = checkConfig(o.Config); err == nil {
		uciAddOpts := client.UCIAddOptions{
			Config: o.Config,
			Type:   o.Type,
		}
		ctx := c.Context()
		rpc := client.GetFromContext(c.Context())
		response, err := rpc.UCI().Add(ctx, uciAddOpts)
		if err != nil {
			return err
		}
		result, err := uciAddOpts.GetResult(response)
		if err != nil {
			return err
		}
		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))
	}
	return err
}

func NewApplyCommand() *cobra.Command {
	o := ApplyOptions{}
	structType := reflect.TypeOf(o)
	numOptions := structType.NumField()
	c := &cobra.Command{
		Use:   "apply",
		Short: "Apply all pending changes.",
		Args:  cobra.MaximumNArgs(numOptions),
		RunE: func(c *cobra.Command, args []string) error {
			return o.Run(c)
		},
	}
	o.BindFlags(c)

	return c
}

type ApplyOptions struct {
	Rollback bool
	Timeout  int
}

func (o *ApplyOptions) BindFlags(c *cobra.Command) {
	c.Flags().BoolVarP(&o.Rollback, "rollback", "r", true, "Undo config if an error is encountered.")
	c.Flags().IntVarP(&o.Timeout, "timeout", "t", 10, "The amount of time to wait before rolling back.")
}

func (o *ApplyOptions) Run(c *cobra.Command) error {
	uciApplyOpts := client.UCIApplyOptions{
		Rollback: uci.ToStringBool(o.Rollback),
		Timeout:  o.Timeout,
	}
	ctx := c.Context()
	rpc := client.GetFromContext(c.Context())
	response, err := rpc.UCI().Apply(ctx, uciApplyOpts)
	if err != nil {
		return err
	}

	output, err := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(output))
	return err
}

func NewChangesCommand() *cobra.Command {
	o := ChangesOptions{}
	structType := reflect.TypeOf(o)
	numOptions := structType.NumField()
	c := &cobra.Command{
		Use:   "changes",
		Short: "List all pending config changes.",
		Args:  cobra.MaximumNArgs(numOptions),
		RunE: func(c *cobra.Command, args []string) error {
			return o.Run(c)
		},
	}
	o.BindFlags(c)

	return c
}

type ChangesOptions struct {
	Config string
}

func (o *ChangesOptions) BindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.Config, "config", "c", "", "Which config to query.")
}

func (o *ChangesOptions) Run(c *cobra.Command) (err error) {
	if err = checkConfig(o.Config); err == nil {
		uciChangesOpts := client.UCIChangesOptions{
			Config: o.Config,
		}
		ctx := c.Context()
		rpc := client.GetFromContext(c.Context())
		response, err := rpc.UCI().Changes(ctx, uciChangesOpts)
		if err != nil {
			return err
		}
		result, err := uciChangesOpts.GetResult(response)
		if err != nil {
			return err
		}
		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))
	}
	return err
}

func NewConfigsCommand() *cobra.Command {
	o := ConfigsOptions{}
	structType := reflect.TypeOf(o)
	numOptions := structType.NumField()
	c := &cobra.Command{
		Use:   "configs",
		Short: "Return a list of all available configs.",
		Args:  cobra.MaximumNArgs(numOptions),
		RunE: func(c *cobra.Command, args []string) error {
			return o.Run(c)
		},
	}

	return c
}

type ConfigsOptions struct{}

func (o *ConfigsOptions) Run(c *cobra.Command) error {
	uciConfigsOpts := client.UCIConfigsOptions{}
	ctx := c.Context()
	rpc := client.GetFromContext(c.Context())
	response, err := rpc.UCI().Configs(ctx, uciConfigsOpts)
	if err != nil {
		return err
	}
	result, err := uciConfigsOpts.GetResult(response)
	if err != nil {
		return err
	}
	output, err := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))
	return err
}

func NewDeleteCommand() *cobra.Command {
	o := DeleteOptions{}
	structType := reflect.TypeOf(o)
	numOptions := structType.NumField()
	c := &cobra.Command{
		Use:   "delete",
		Short: "Delete a config value.",
		Args:  cobra.MaximumNArgs(numOptions),
		RunE: func(c *cobra.Command, args []string) error {
			return o.Run(c)
		},
	}
	o.BindFlags(c)

	return c
}

type DeleteOptions struct {
	Config  string
	Section string
	Type    string
	Option  string
}

func (o *DeleteOptions) BindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.Config, "config", "c", "", "Which config to query.")
	c.Flags().StringVarP(&o.Section, "section", "s", "", "The section of the config.")
	c.Flags().StringVarP(&o.Type, "type", "t", "", "The type of the config section.")
	c.Flags().StringVarP(&o.Option, "option", "o", "", "A single option within a config section.")
	c.MarkFlagRequired("config")
}

func (o *DeleteOptions) Run(c *cobra.Command) (err error) {
	if err = checkConfig(o.Config); err == nil {
		uciDeleteOpts := client.UCIDeleteOptions{
			Config:  o.Config,
			Section: o.Section,
			Type:    o.Type,
			Option:  o.Option,
		}
		ctx := c.Context()
		rpc := client.GetFromContext(c.Context())
		response, err := rpc.UCI().Delete(ctx, uciDeleteOpts)
		if err != nil {
			return err
		}

		output, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(output))
	}
	return err
}

func NewGetCommand() *cobra.Command {
	o := GetOptions{}
	structType := reflect.TypeOf(o)
	numOptions := structType.NumField()
	c := &cobra.Command{
		Use:   "get",
		Short: "Get a config value.",
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

func (o *GetOptions) Run(c *cobra.Command) (err error) {
	if err = checkConfig(o.Config); err == nil {
		uciGetOpts := client.UCIGetOptions{
			Config:  o.Config,
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
		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))
	}
	return err
}

func NewRevertCommand() *cobra.Command {
	o := RevertOptions{}
	structType := reflect.TypeOf(o)
	numOptions := structType.NumField()
	c := &cobra.Command{
		Use:   "revert",
		Short: "Clear pending config changes.",
		Args:  cobra.MaximumNArgs(numOptions),
		RunE: func(c *cobra.Command, args []string) error {
			return o.Run(c)
		},
	}
	o.BindFlags(c)

	return c
}

type RevertOptions struct {
	Config string
}

func (o *RevertOptions) BindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.Config, "config", "c", "", "Which config to query.")
	c.MarkFlagRequired("config")
}

func (o *RevertOptions) Run(c *cobra.Command) (err error) {
	if err = checkConfig(o.Config); err == nil {
		uciRevertOpts := client.UCIRevertOptions{
			Config: o.Config,
		}
		ctx := c.Context()
		rpc := client.GetFromContext(c.Context())
		response, err := rpc.UCI().Revert(ctx, uciRevertOpts)
		if err != nil {
			return err
		}
		output, err := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(output))
	}
	return err
}

func NewSetCommand() *cobra.Command {
	o := SetOptions{}
	structType := reflect.TypeOf(o)
	numOptions := structType.NumField()
	c := &cobra.Command{
		Use:   "set",
		Short: "Set config values.",
		Args:  cobra.MaximumNArgs(numOptions),
		RunE: func(c *cobra.Command, args []string) error {
			return o.Run(c)
		},
	}
	o.BindFlags(c)

	return c
}

type SetOptions struct {
	Config  string
	Section string
	Type    string
	Values  string // '{"enabled": "1"}'
}

func (o *SetOptions) BindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.Config, "config", "c", "", "Which config to query.")
	c.Flags().StringVarP(&o.Section, "section", "s", "", "The section of the config.")
	c.Flags().StringVarP(&o.Type, "type", "t", "", "The type of section in the config.")
	c.Flags().StringVarP(&o.Values, "values", "v", "", "The option: value pairs to set, must be JSON-formatted.")
	c.MarkFlagRequired("config")
	c.MarkFlagRequired("section")
	c.MarkFlagRequired("type")
	c.MarkFlagRequired("values")
}

func (o *SetOptions) Run(c *cobra.Command) (err error) {
	if err = checkConfig(o.Config); err == nil {
		uciSetOpts := client.UCISetOptions{}
		switch o.Type {
		case string(dhcp.Boot):
			uciSetOpts = unmarshalCLIValues[dhcp.BootSection](o)
		case string(dhcp.CircuitID):
			uciSetOpts = unmarshalCLIValues[dhcp.CircuitIDSection](o)
		case string(dhcp.DHCP):
			uciSetOpts = unmarshalCLIValues[dhcp.DHCPSection](o)
		case string(dhcp.Dnsmasq):
			uciSetOpts = unmarshalCLIValues[dhcp.DnsmasqSection](o)
		case string(dhcp.Host):
			uciSetOpts = unmarshalCLIValues[dhcp.HostSection](o)
		case string(dhcp.HostRecord):
			uciSetOpts = unmarshalCLIValues[dhcp.HostRecordSection](o)
		case string(dhcp.MAC):
			uciSetOpts = unmarshalCLIValues[dhcp.MACSection](o)
		case string(dhcp.Odhcpd):
			uciSetOpts = unmarshalCLIValues[dhcp.OdhcpdSection](o)
		case string(dhcp.Relay):
			uciSetOpts = unmarshalCLIValues[dhcp.RelaySection](o)
		case string(dhcp.RemoteID):
			uciSetOpts = unmarshalCLIValues[dhcp.RemoteIDSection](o)
		case string(dhcp.SubscrID):
			uciSetOpts = unmarshalCLIValues[dhcp.SubscrIDSection](o)
		case string(dhcp.Tag):
			uciSetOpts = unmarshalCLIValues[dhcp.TagSection](o)
		case string(dhcp.UserClass):
			uciSetOpts = unmarshalCLIValues[dhcp.UserClassSection](o)
		case string(dhcp.VendorClass):
			uciSetOpts = unmarshalCLIValues[dhcp.VendorClassSection](o)
		case string(firewall.Defaults):
			uciSetOpts = unmarshalCLIValues[firewall.DefaultsSectionOptions](o)
		case string(firewall.Forwarding):
			uciSetOpts = unmarshalCLIValues[firewall.ForwardingSectionOptions](o)
		case string(firewall.Redirect):
			uciSetOpts = unmarshalCLIValues[firewall.RedirectSectionOptions](o)
		case string(firewall.Rule):
			uciSetOpts = unmarshalCLIValues[firewall.RuleSectionOptions](o)
		case string(firewall.Zone):
			uciSetOpts = unmarshalCLIValues[firewall.ZoneSectionOptions](o)
		}

		ctx := c.Context()
		rpc := client.GetFromContext(c.Context())
		response, err := rpc.UCI().Set(ctx, uciSetOpts)
		if err != nil {
			return err
		}

		output, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(output))
	}
	return err
}
