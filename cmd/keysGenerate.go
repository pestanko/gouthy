/*command
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/cmd/cmd_utils"

	"github.com/spf13/cobra"
)

// keysGenerateCmd represents the keysGenerate command
var keysGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new RSA Key",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers apps.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd_utils.BindAppContext(generateNewKey, cmd, args)
	},
}

func generateNewKey(ctx context.Context, app *core.GouthyApp, cmd *cobra.Command, args []string) error {
	err := app.Facades.Keys.GenerateNewJwk(ctx)
	return err
}

func init() {
	keysCmd.AddCommand(keysGenerateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keysGenerateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// keysGenerateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
