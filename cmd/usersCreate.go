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
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/pestanko/gouthy/cmd/cmd_utils"
	"github.com/spf13/cobra"
)

var user users.NewUserDTO
var secret string

// usersCreateCmd represents the usersCreate command
var usersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd_utils.BindAppContext(createNewUser, cmd, args)
	},
}

func init() {
	usersCmd.AddCommand(usersCreateCmd)
	usersCreateCmd.PersistentFlags().StringVarP(&user.Email, "email", "e", "", "User Email")
	usersCreateCmd.PersistentFlags().StringVarP(&user.Name, "name", "n", "", "User Full name")
	usersCreateCmd.PersistentFlags().StringVarP(&user.Username, "username", "u", "", "User's username")
	usersCreateCmd.PersistentFlags().StringVarP(&user.Password, "password", "p", "", "User's password (not recommended)")
	usersCreateCmd.PersistentFlags().StringVarP(&secret, "secret", "S", "", "User's secret")
}

func createNewUser(ctx context.Context, app *infra.GouthyApp, cmd *cobra.Command, args []string) error {
	newUser, err := app.Facades.Users.Create(ctx, &user)

	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(&newUser, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("New user: %v \n", data)
	return nil
}
