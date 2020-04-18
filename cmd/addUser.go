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
	"encoding/json"
	"fmt"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/services"
	"github.com/pestanko/gouthy/cmd/helpers"
	"github.com/spf13/cobra"
)

var user services.NewUser
var secret string

// addUserCmd represents the addUser command
var addUserCmd = &cobra.Command{
	Use:   "add-user",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		helpers.BindAppContext(createNewUser, cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(addUserCmd)
	addUserCmd.PersistentFlags().StringVarP(&user.Email, "email", "e", "", "User Email")
	addUserCmd.PersistentFlags().StringVarP(&user.Name, "name", "n", "", "User Full name")
	addUserCmd.PersistentFlags().StringVarP(&user.Username, "username", "u", "", "User's username")
	addUserCmd.PersistentFlags().StringVarP(&user.Password, "password", "p", "", "User's password (not recommended)")
	addUserCmd.PersistentFlags().StringVarP(&secret, "secret", "S", "", "User's secret")
}

func createNewUser(app *core.GouthyApp, cmd *cobra.Command, args []string) error {

	newUser, err := app.Services.Users.Create(&user)

	if err != nil {
		return err
	}

	cvtUser := services.ConvertModelsToUserDTO(&newUser)

	data, err := json.MarshalIndent(&cvtUser, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("New user: %v \n", data)
	return nil
}
