package core

import (
	"context"
	"encoding/json"
	"github.com/pestanko/gouthy/app/apps"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/users"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ImportSchema struct {
	Users []users.CreateDTO `json:"users" yaml:"users"`
	Apps  []apps.CreateDTO  `json:"apps" yaml:"apps"`
}

func NewDataImporter(app *GouthyApp) *DataImporter {
	return &DataImporter{
		app: app,
	}
}

type DataImporter struct {
	app *GouthyApp
}

func (imp *DataImporter) ImportFromFile(ctx context.Context, path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	all, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	schema := ImportSchema{}
	ext := filepath.Ext(path)
	switch ext {
	case "yml", "yaml":
		if err := yaml.Unmarshal(all, &schema); err != nil {
			return err
		}
	case "json":
		if err := json.Unmarshal(all, &schema); err != nil {
			return err
		}
	default:
		shared.GetLogger(ctx).WithFields(logrus.Fields{
			"file": path,
			"ext":  ext,
		}).Error("Unsupported file format")
		return nil
	}

	return imp.importSchema(ctx, &schema)
}

func (imp *DataImporter) importSchema(ctx context.Context, schema *ImportSchema) error {
	// import users
	if err := imp.importUsers(ctx, schema.Users); err != nil {
		return err
	}

	// import apps
	if err := imp.importApps(ctx, schema.Apps); err != nil {
		return err
	}

	return nil
}

func (imp *DataImporter) importUsers(ctx context.Context, dtos []users.CreateDTO) error {
	shared.GetLogger(ctx).Info("Importing users")
	for _, dto := range dtos {
		fields := shared.GetLogger(ctx).WithFields(dto.LogFields())
		fields.Debug("Importing user")
		user, err := imp.app.Facades.Users.Create(ctx, &dto)
		if err != nil {
			fields.Error("Unable to import user")
		}
		fields.WithField("user_id", user.ID.String()).Info("Imported user")
	}

	return nil
}

func (imp *DataImporter) importApps(ctx context.Context, dtos []apps.CreateDTO) error {
	shared.GetLogger(ctx).Info("Importing apps")
	for _, dto := range dtos {
		fields := shared.GetLogger(ctx).WithFields(dto.LogFields())
		fields.Debug("Importing apps")
		app, err := imp.app.Facades.Apps.Create(ctx, &dto)
		if err != nil {
			fields.Error("Unable to import apps")
		}
		fields.WithField("app_id", app.ID.String()).Info("Imported app")
	}

	return nil
}

func (imp *DataImporter) ImportFromFiles(ctx context.Context, files []string) error {
	for _, file := range files {
		if err := imp.ImportFromFile(ctx, file); err != nil {
			shared.GetLogger(ctx).WithError(err).WithField("file", file).Error("Unable to import from file")
		}
	}
	return nil
}
