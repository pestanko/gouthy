package auth

import (
	"github.com/pestanko/gouthy/app/shared"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	StepLoginPassword = "login_password"
	StepLoginTotp     = "login_totp"
	StepLoginSecret   = "login_secret"
	StepFindUser      = "find_user"
)

func NewLoginStep(name string, state LoginStateStatus) LoginStep {
	return &loginStepImpl{name: name, state: state}
}

type LoginStep interface {
	Name() string
	State() LoginStateStatus
	IsSuccess() bool
	IsFail() bool
	IsError() bool
}

type loginStepImpl struct {
	name  string
	state LoginStateStatus
}

func (step *loginStepImpl) IsError() bool {
	return step.state == Error
}

func (step *loginStepImpl) State() LoginStateStatus {
	return step.state
}

func (step *loginStepImpl) IsSuccess() bool {
	return step.state == Success
}

func (step *loginStepImpl) IsFail() bool {
	return !step.IsSuccess()
}
func (step *loginStepImpl) Name() string {
	return step.name
}

type LoginStateStatus string

const (
	Success LoginStateStatus = "success"
	Failed  LoginStateStatus = "failed"
	Error   LoginStateStatus = "error"
)

type LoginState interface {
	shared.LoggingIdentity
	UserID() string
	ID() string
	Steps() *[]LoginStep
	AddStep(loginStep LoginStep) LoginState
	IsSuccess() bool
	IsFail() bool
	IsError() bool
	State() LoginStateStatus
	IsNotOk() bool
	IsOk() bool
}

type loginStateImpl struct {
	loginSteps []LoginStep
	userId     string
	id         string
}

func (state *loginStateImpl) LogFields() log.Fields {
	return log.Fields{
		"state": state.State(),
		"user_id": state.UserID(),
		"login_state_id": state.ID(),
	}
}

func (state *loginStateImpl) IsNotOk() bool {
	return !state.IsOk()
}

func (state *loginStateImpl) ID() string {
	return state.id
}

func (state *loginStateImpl) UserID() string {
	return state.userId
}

func (state *loginStateImpl) IsSuccess() bool {
	return state.State() == Success
}

func (state *loginStateImpl) IsFail() bool {
	return state.State() == Failed
}

func (state *loginStateImpl) IsError() bool {
	return state.State() == Error
}

func (state *loginStateImpl) State() LoginStateStatus {
	for _, item := range state.loginSteps {
		if item.State() != Success {
			return item.State()
		}
	}
	return Success
}

func (state *loginStateImpl) GetUnsuccessfulStep() LoginStep {
	for _, item := range state.loginSteps {
		if item.State() != Success {
			return item
		}
	}
	return nil
}

func (state *loginStateImpl) GetStepByName(name string) LoginStep {
	for _, item := range state.loginSteps {
		if item.Name() == name {
			return item
		}
	}
	return nil
}

func (state *loginStateImpl) Steps() *[]LoginStep {
	return &state.loginSteps
}

func (state *loginStateImpl) AddStep(loginStep LoginStep) LoginState {
	state.loginSteps = append(state.loginSteps, loginStep)
	return state
}

func (state *loginStateImpl) IsOk() bool {
	return state.IsSuccess()
}

func NewLoginState(userId uuid.UUID) LoginState {
	return &loginStateImpl{
		loginSteps: []LoginStep{},
		userId:     userId.String(),
		id:         uuid.NewV4().String(),
	}
}
