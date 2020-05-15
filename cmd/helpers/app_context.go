package helpers

import (
	"fmt"
	"github.com/pestanko/gouthy/app/infra"
	"github.com/spf13/cobra"
	"os"
)

func BindAppContext(fn func(app *infra.GouthyApp, cmd *cobra.Command, args []string) error,  cmd *cobra.Command, args []string) {
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

	if err = fn(&app, cmd, args); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

}