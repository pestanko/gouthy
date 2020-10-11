package auth

import (
	"context"
	"fmt"
	"github.com/pestanko/gouthy/app/apps"
	"github.com/pestanko/gouthy/app/shared"
	"github.com/pestanko/gouthy/app/shared/utils"
	log "github.com/sirupsen/logrus"
)

type AuthorizationResult struct {
}

type OAuth2AuthorizationService interface {
	ProcessRequest(ctx context.Context, request *OAuth2AuthRequest) (AuthorizationResult, error)
}

func NewOAuth2AuthorizationService(appFindService apps.FindService) OAuth2AuthorizationService {
	return &oAuth2AuthorizationServiceImpl{appFindService: appFindService}
}

type oAuth2AuthorizationServiceImpl struct {
	appFindService apps.FindService
}

func (s *oAuth2AuthorizationServiceImpl) ProcessRequest(ctx context.Context, request *OAuth2AuthRequest) (AuthorizationResult, error) {
	if err := s.checkRequest(ctx, request); err != nil {
		return AuthorizationResult{}, err
	}

	if err := s.checkApplication(ctx, request); err != nil {

	}

	return AuthorizationResult{}, nil
}

func (s *oAuth2AuthorizationServiceImpl) checkRequest(ctx context.Context, request *OAuth2AuthRequest) error {
	if request.ClientId == "" {
		return s.failParamMissing(ctx, request, "client_id")
	}

	if request.State == "" {
		return s.failParamMissing(ctx, request, "state")
	}

	if request.ResponseType == "" {
		return s.failParamMissing(ctx, request, "response_type")
	}

	if request.RedirectUri == "" {
		return s.failParamMissing(ctx, request, "redirect_uri")
	}

	if request.PKCEMethod != "" || request.PKCEChallenge != "" {
		if request.PKCEMethod == "" {
			return s.failParamMissing(ctx, request, "pkce_code_method")
		}
		if request.PKCEChallenge == "" {
			return s.failParamMissing(ctx, request, "pkce_code_challenge")
		}
	}

	return nil
}

func (s *oAuth2AuthorizationServiceImpl) checkApplication(ctx context.Context, request *OAuth2AuthRequest) error {
	app, err := s.appFindService.FindOne(ctx, apps.FindQuery{ClientId: request.ClientId})
	if err != nil {
		return err
	}

	if app == nil {
		return NewOAuth2AuthorizationRequestError("Application not found for provided client_id").WithDetail(shared.ErrDetail{
			"client_id": request.ClientId,
			"status":    shared.ErrStatusInvalidRequest,
		})
	}

	if !app.IsActive() {
		return apps.NewApplicationError("Application is not active").WithDetail(shared.ErrDetail{
			"client_id": request.ClientId,
			"status":    shared.ErrStatusInvalidRequest,
		})
	}

	diff := utils.StringSliceDifference(request.Scopes, app.Scopes())
	if len(diff) != 0 {
		return NewOAuth2AuthorizationRequestError("Request contains unsupported scopes").WithDetail(shared.ErrDetail{
			"additional_scopes": diff,
			"status":            shared.ErrStatusInvalidRequest,
		})
	}

	return nil
}

func (s *oAuth2AuthorizationServiceImpl) failParamMissing(ctx context.Context, request *OAuth2AuthRequest, param string) error {
	shared.GetLogger(ctx).WithFields(log.Fields{
		"param":     param,
		"client_id": request.ClientId,
	}).Error("missing required parameter")
	err := NewOAuth2AuthorizationRequestError(fmt.Sprintf("Required parameter is missing: %v", param))
	return err
}
