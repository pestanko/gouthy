package cmd_utils

import (
	"context"
	"fmt"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/spf13/cobra"
	"os"
)

func BindAppContext(fn func(ctx context.Context, app *infra.GouthyApp, cmd *cobra.Command, args []string) error, cmd *cobra.Command, args []string) {
	var err error
	config, err := infra.GetAppConfig()
	checkError(err)


	db, err := infra.GetDBConnection(&config)
	checkError(err)

	defer db.Close()

	app, err := infra.GetApplication(&config, db)
	checkError(err)

	ctx := shared.NewOperationContext()
	ctx = context.WithValue(ctx, "client_id", "admin_console")

	checkError(fn(ctx, &app, cmd, args))
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("CLI error: %v", err)
		os.Exit(1)
	}
}