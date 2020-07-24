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
	"fmt"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/pestanko/gouthy/app/web"
	"github.com/pestanko/gouthy/cmd/cmd_utils"

	"github.com/spf13/cobra"
)

// routesCmd represents the routes command
var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd_utils.BindAppContext(runRoutes, cmd, args)
	},
}

func runRoutes(ctx context.Context, app *infra.GouthyApp, cmd *cobra.Command, args []string) error {
	webServer := web.CreateWebServer(app)

	web.RegisterRoutes(webServer)

	fmt.Println("Routes:")
	for _, route := range webServer.Router.Routes() {
		fmt.Printf("%v - %v\n", route.Method, route.Path)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(routesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// routesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// routesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
