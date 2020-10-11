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
	"github.com/pestanko/gouthy/app/apps"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/cmd/cmd_utils"

	"github.com/spf13/cobra"
)

var appDTO apps.CreateDTO
var appSecret string

// appsCreateCmd represents the appsCreate command
var appsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new applicationModel",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers apps.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd_utils.BindAppContext(createNewApp, cmd, args)
	},
}

func createNewApp(ctx context.Context, app *core.GouthyApp, cmd *cobra.Command, args []string) error {
	newApp, err := app.Facades.Apps.Create(ctx, &appDTO)

	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(&newApp, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

func init() {
	appsCmd.AddCommand(appsCreateCmd)

	appsCreateCmd.PersistentFlags().StringVarP(&appDTO.Description, "desc", "D", "", "applicationModel Description")
	appsCreateCmd.PersistentFlags().StringVarP(&appDTO.Name, "name", "n", "", "applicationModel name")
	appsCreateCmd.PersistentFlags().StringVarP(&appDTO.Codename, "codename", "c", "", "applicationModel's codename")
	appsCreateCmd.PersistentFlags().StringVarP(&appDTO.ClientId, "client-id", "C", "", "applicationModel's clientId")
	appsCreateCmd.PersistentFlags().StringVarP(&appDTO.Type, "type", "T", "", "applicationModel's appSecret")

	appsCreateCmd.PersistentFlags().StringVarP(&appSecret, "appSecret", "S", "", "applicationModel's appSecret")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appsCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appsCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
