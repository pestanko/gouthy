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
	"bufio"
	"context"
	"fmt"
	"github.com/pestanko/gouthy/app/domain/auth"
	"github.com/pestanko/gouthy/app/domain/jwtlib"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/pestanko/gouthy/cmd/cmd_utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// loginPasswordCmd represents the loginPassword command
var loginPasswordCmd = &cobra.Command{
	Use:   "password",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd_utils.BindAppContext(executeLoginPassword, cmd, args)
	},
}

func executeLoginPassword(ctx context.Context, app *infra.GouthyApp, cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if err != nil {
		return err
	}

	password, err := cmd_utils.RequestPassword()
	if err != nil {
		return err
	}

	user, err := app.Facades.Users.GetByUsername(ctx, username)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("unable to find user: %s", username)
	}

	loginState := auth.NewLoginState(user.ID)
	loginState, err = app.Facades.Auth.LoginUsernamePassword(ctx, loginState, auth.PasswordLoginDTO{
		Username: username,
		Password: password,
	})

	if loginState == nil || loginState.IsFail() {
		return fmt.Errorf("fogin failed for user %s", username)
	}

	application, err := app.Facades.Apps.GetByAnyId(ctx, AdminConsoleApp)
	if err != nil {
		return err
	}
	if application == nil {
		return fmt.Errorf("unable to find application: %s", AdminConsoleApp)
	}

	params := jwtlib.TokenCreateParams{
		User:   user,
		App:    application,
		Scopes: []string{"console_login"},
	}
	accessToken, err := app.Facades.Jwt.CreateSignedAccessToken(ctx, params)
	if err != nil {
		return err
	}

	fmt.Println(accessToken.Signature)

	return nil
}

func init() {
	loginCmd.AddCommand(loginPasswordCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginPasswordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginPasswordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
