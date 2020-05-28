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
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	db, err := infra.GetDBConnection(&config)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	defer db.Close()

	app, err := infra.GetApplication(&config, db)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	ctx := shared.NewOperationContext()
	ctx = context.WithValue(ctx, "client_id", "admin_console")

	if err = fn(ctx, &app, cmd, args); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
