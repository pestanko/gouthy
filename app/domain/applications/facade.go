package applications

import (
	"context"
	"github.com/pestanko/gouthy/app/shared"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type ListParams struct {
}

type Facade interface {
	Create(ctx context.Context, newApp *CreateDTO) (*ApplicationDTO, error)
	Update(ctx context.Context, userId uuid.UUID, newUser *UpdateDTO) (*ApplicationDTO, error)
	Delete(ctx context.Context, userId uuid.UUID) error
	List(ctx context.Context, listParams ListParams) ([]ListApplicationDTO, error)
	Get(ctx context.Context, appId uuid.UUID) (*ApplicationDTO, error)
	GetByCodename(ctx context.Context, appId string) (*ApplicationDTO, error)
	GetByAnyId(ctx context.Context, sid string) (*ApplicationDTO, error)
}

type facadeImpl struct {
	apps    Repository
	secrets SecretsRepository
}

func (f *facadeImpl) Create(ctx context.Context, newApp *CreateDTO) (*ApplicationDTO, error) {

	clientId := newApp.ClientId
	if clientId == "" {
		clientId = uuid.NewV4().String()
	}

	var app = &Application{
		Codename:    newApp.Codename,
		Name:        newApp.Name,
		Description: newApp.Description,
		Type:        newApp.Type,
		ClientId:    clientId,
	}

	if err := f.apps.Create(ctx, app); err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"codename":  app.Codename,
			"client_id": app.ClientId,
		}).Error("Unable to create a application")
		return nil, err
	}

	shared.GetLogger(ctx).WithFields(log.Fields{
		"application_id": app.ID,
		"client_id":      app.ClientId,
		"codename":       app.Codename,
	}).Info("Creating a new application")

	return ConvertModelToDTO(app), nil
}

func (f *facadeImpl) Update(ctx context.Context, appId uuid.UUID, newApp *UpdateDTO) (*ApplicationDTO, error) {
	var app = &Application{
		ID:          appId,
		Codename:    newApp.Codename,
		Name:        newApp.Name,
		Description: newApp.Description,
	}

	if err := f.apps.Update(ctx, app); err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"codename":  app.Codename,
			"client_id": app.ClientId,
		}).Error("Unable to create a application")
		return nil, err
	}

	shared.GetLogger(ctx).WithFields(log.Fields{
		"application_id": app.ID,
		"client_id":      app.ClientId,
		"codename":       app.Codename,
	}).Info("Creating a new application")

	return ConvertModelToDTO(app), nil
}

func (f *facadeImpl) Delete(ctx context.Context, userId uuid.UUID) error {
	var app, err = f.apps.FindByID(ctx, userId)
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"app_id":   app.ID,
			"codename": app.Codename,
		}).Error("Unable to delete an app")
		return err
	}

	return f.apps.Delete(ctx, app)
}

func (f *facadeImpl) List(ctx context.Context, listParams ListParams) ([]ListApplicationDTO, error) {
	list, err := f.apps.List(ctx)
	if err != nil {
		return []ListApplicationDTO{}, err
	}

	return ConvertModelsToList(list), err
}

func (f *facadeImpl) Get(ctx context.Context, appId uuid.UUID) (*ApplicationDTO, error) {
	var app, err = f.apps.FindByID(ctx, appId)
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"app_id": appId,
		}).Error("Unable to get an app")
		return nil, err
	}

	return ConvertModelToDTO(app), nil
}

func (f *facadeImpl) GetByCodename(ctx context.Context, appId string) (*ApplicationDTO, error) {
	var app, err = f.apps.FindByCodename(ctx, appId)
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"codename": appId,
		}).Error("Unable to get an app")
		return nil, err
	}

	return ConvertModelToDTO(app), nil
}

func (f *facadeImpl) GetByAnyId(ctx context.Context, sid string) (*ApplicationDTO, error) {
	var uid, err = uuid.FromString(sid)
	if err == nil {
		return f.Get(ctx, uid)
	}

	return f.GetByCodename(ctx, sid)
}

func NewApplicationsFacade(apps Repository, secrets SecretsRepository) Facade {
	return &facadeImpl{apps: apps, secrets: secrets}
}
