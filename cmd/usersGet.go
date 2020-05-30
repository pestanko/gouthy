/*
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
	"encoding/json"
	"fmt"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/pestanko/gouthy/cmd/cmd_utils"

	"github.com/spf13/cobra"
)

// usersGetCmd represents the usersGet command
var usersGetCmd = &cobra.Command{
	Use:   "usersGet",
	Short: "Get a user info",
	Long: `Gets a user info in JSON format. It is possible to provide multiple user IDs.
They can be either uuid or username.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd_utils.BindAppContext(getUser, cmd, args)
	},
}

func getUser(ctx context.Context, app *infra.GouthyApp, cmd *cobra.Command, args []string) error {
	for _, username := range args {
		user, err := app.Facades.Users.GetByAnyId(ctx, username)
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(&user, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))

	}
	return nil
}

func init() {
	usersCmd.AddCommand(usersGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
