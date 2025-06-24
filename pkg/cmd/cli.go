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
