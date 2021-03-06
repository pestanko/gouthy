package apps

import (
	"context"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/shared/repos"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type ListParams struct {
	Offset int
	Limit  int
}

type Facade interface {
	Create(ctx context.Context, newApp *CreateDTO) (*Application, error)
	Update(ctx context.Context, userId uuid.UUID, newUser *UpdateDTO) (*Application, error)
	Delete(ctx context.Context, userId uuid.UUID) error
	List(ctx context.Context, listParams ListParams) ([]Application, error)
	Get(ctx context.Context, appId uuid.UUID) (*Application, error)
	GetByCodename(ctx context.Context, appId string) (*Application, error)
	GetByAnyId(ctx context.Context, sid string) (*Application, error)
	GetByClientId(ctx context.Context, sid string) (*Application, error)
}


func NewApplicationsFacade(apps Repository, secrets SecretsRepository, findService FindService) Facade {
	return &facadeImpl{apps: apps, secrets: secrets, findService: findService}
}


type facadeImpl struct {
	apps        Repository
	secrets     SecretsRepository
	findService FindService
}

func (f *facadeImpl) GetByClientId(ctx context.Context, clientId string) (*Application, error) {
	var app, err = f.apps.QueryOne(ctx, FindQuery{ClientId: clientId})
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"client_id": clientId,
		}).Error("Unable to get an app")
		return nil, err
	}

	return app, nil
}

func (f *facadeImpl) Create(ctx context.Context, newApp *CreateDTO) (*Application, error) {

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

	return app, nil
}

func (f *facadeImpl) Update(ctx context.Context, appId uuid.UUID, newApp *UpdateDTO) (*Application, error) {
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

	return app, nil
}

func (f *facadeImpl) Delete(ctx context.Context, appId uuid.UUID) error {
	var app, err = f.apps.QueryOne(ctx, FindQuery{Id: appId})
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"application_id": app.ID,
			"codename":       app.Codename,
		}).Error("Unable to delete an app")
		return err
	}

	return f.apps.Delete(ctx, app)
}

func (f *facadeImpl) List(ctx context.Context, params ListParams) ([]Application, error) {
	list, err := f.apps.Query(ctx, FindQuery{
		PaginationQuery: repos.NewPaginationQuery(params.Limit, params.Offset),
	})
	if err != nil {
		return []Application{}, err
	}

	return list, err
}

func (f *facadeImpl) Get(ctx context.Context, appId uuid.UUID) (*Application, error) {
	var app, err = f.apps.QueryOne(ctx, FindQuery{Id: appId})
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"application_id": appId,
		}).Error("Unable to get an app")
		return nil, err
	}

	return app, nil
}

func (f *facadeImpl) GetByCodename(ctx context.Context, codename string) (*Application, error) {
	var app, err = f.apps.QueryOne(ctx, FindQuery{Codename: codename})
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"codename": codename,
		}).Error("Unable to get an app")
		return nil, err
	}

	return app, nil
}

func (f *facadeImpl) GetByAnyId(ctx context.Context, sid string) (*Application, error) {
	one, err := f.findService.FindOne(ctx, FindQuery{AnyId: sid})
	if err != nil {
		shared.GetLogger(ctx).WithError(err).WithFields(log.Fields{
			"any_id": sid,
		}).Error("Unable to get an app")
		return nil, err
	}
	return one, err
}
