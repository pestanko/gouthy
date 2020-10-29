package data

import (
	"context"
	"encoding/json"
	"github.com/pestanko/gouthy/app/core"
	"github.com/pestanko/gouthy/app/domain/apps"
	"github.com/pestanko/gouthy/app/domain/users"
	"github.com/pestanko/gouthy/app/shared"
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

func NewImporter(app *core.GouthyApp) *Importer {
	return &Importer{
		app: app,
	}
}

type Importer struct {
	app *core.GouthyApp
}

func (imp *Importer) ImportFromFile(ctx context.Context, path string) error {
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
	case ".yml", ".yaml":
		if err := yaml.Unmarshal(all, &schema); err != nil {
			return err
		}
	case ".json":
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

func (imp *Importer) importSchema(ctx context.Context, schema *ImportSchema) error {
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

func (imp *Importer) importUsers(ctx context.Context, dtos []users.CreateDTO) error {
	shared.GetLogger(ctx).Info("Importing users")
	for _, dto := range dtos {
		fields := shared.GetLogger(ctx).WithFields(dto.LogFields())

		// check whether user exists, if so - continue
		found, err := imp.app.Facades.Users.GetByUsername(ctx, dto.Username)
		if err != nil || found != nil {
			fields.Warn("User already exists, skipping")
			continue
		}

		fields.Debug("Importing user")
		user, err := imp.app.Facades.Users.Create(ctx, &dto)
		if err != nil {
			fields.Error("Unable to import user")
			continue
		}
		fields.WithField("user_id", user.ID.String()).Info("Imported user")
	}

	return nil
}

func (imp *Importer) importApps(ctx context.Context, dtos []apps.CreateDTO) error {
	shared.GetLogger(ctx).Info("Importing apps")
	for _, dto := range dtos {
		fields := shared.GetLogger(ctx).WithFields(dto.LogFields())

		// check whether user exists, if so - continue
		found, err := imp.app.Facades.Apps.GetByCodename(ctx, dto.Codename)
		if err != nil || found == nil {
			// TODO: Override option
			fields.Warn("App already exists, skipping")
			continue
		}

		fields.Debug("Importing app")
		app, err := imp.app.Facades.Apps.Create(ctx, &dto)
		if err != nil {
			fields.Error("Unable to import app")
		}
		fields.WithField("app_id", app.ID.String()).Info("Imported app")
	}

	return nil
}

func (imp *Importer) ImportFromFiles(ctx context.Context, files []string) error {
	for _, file := range files {
		logFields := shared.GetLogger(ctx).WithField("file", file)
		logFields.Info("Importing from file started")
		if err := imp.ImportFromFile(ctx, file); err != nil {
			logFields.WithError(err).Error("Unable to import from file")
		}
		logFields.Debug("Importing from file ended")
	}
	return nil
}
