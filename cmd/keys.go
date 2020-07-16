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
	"github.com/jedib0t/go-pretty/table"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/pestanko/gouthy/cmd/cmd_utils"

	"github.com/spf13/cobra"
)

// keysCmd represents the keys command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage JWKs or list when no additional action has been provided",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers apps.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd_utils.BindAppContext(listKeys, cmd, args)
	},
}

func listKeys(ctx context.Context, app *infra.GouthyApp, cmd *cobra.Command, args []string) error {
	list, err := app.Facades.Keys.ListJwks(ctx)
	if err != nil {
		return err
	}
	// https://github.com/jedib0t/go-pretty/blob/master/cmd/demo-table/demo.go
	tw := table.NewWriter()
	tw.SetTitle("Gouthy users")
	tw.SetIndexColumn(1)
	tw.AppendHeader(table.Row{"#", "ID", "Algorithm"})

	for i, entity := range list {
		tw.AppendRow(table.Row{i, entity.KeyId(), entity.Algorithm()})
	}

	fmt.Println(tw.Render())

	return nil
}

func init() {
	rootCmd.AddCommand(keysCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keysCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// keysCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
