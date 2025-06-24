/*
Copyright 2025 Daimonas Labs.

Licensed under the GNU General Public License, Version 3 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.gnu.org/licenses/gpl-3.0.en.html

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
