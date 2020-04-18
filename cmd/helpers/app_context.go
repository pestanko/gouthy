package helpers

import (
	"fmt"
	"github.com/pestanko/gouthy/app/core"
	"github.com/spf13/cobra"
	"os"
)

func BindAppContext(fn func(app *core.GouthyApp, cmd *cobra.Command, args []string) error,  cmd *cobra.Command, args []string) {
	var err error
	config, err := core.GetAppConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	db, err := core.GetDBConnection(&config)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	defer db.Close()

	app, err := core.GetApplication(&config, db)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	if err = fn(&app, cmd, args); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

}