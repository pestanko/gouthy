package cmd_utils

import (
	"context"
	"fmt"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/tools/data"
	"github.com/spf13/cobra"
	"os"
)

func BindAppContext(fn func(ctx context.Context, app *core.GouthyApp, cmd *cobra.Command, args []string) error, cmd *cobra.Command, args []string) {
	var err error
	config, err := shared.GetAppConfig()
	checkError(err)
	ctx := shared.NewContextWithConfiguration(&config)
	ctx = context.WithValue(ctx, "client_id", "admin_console")

	db, err := shared.GetDBConnection(&config)
	checkError(err)
	defer db.Close()

	defaultDBConfig := config.DB.GetDefault()
	if config.DB.AutoMigrate && defaultDBConfig.AutoMigrate {
		shared.GetLogger(ctx).Info("Starting migration")
		err = migrate(ctx, db)
		checkError(err)
	}

	app, err := core.GetApplication(&config, db)
	checkError(err)

	if len(defaultDBConfig.DataImport) > 0 {
		_ = data.NewImporter(app).ImportFromFiles(ctx, defaultDBConfig.DataImport)
	}

	checkError(fn(ctx, app, cmd, args))
}

func migrate(ctx context.Context, db shared.DBConnection) error {
	gormDB := shared.DBConnectionIntoGorm(db)
	return gormDB.Migrator().AutoMigrate(
		&users.User{},
		&apps.Application{},
		&apps.SecretModel{},
	)
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("CLI error: %v\n", err)
		os.Exit(1)
	}
}
