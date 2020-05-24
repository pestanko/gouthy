package auth

import uuid "github.com/satori/go.uuid"

type LoginStep interface {
	Name() string
	State() LoginStateStatus
	IsSuccess() bool
	IsFail() bool
}

type loginStepImpl struct {
	name  string
	state LoginStateStatus
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

func NewLoginStep(name string, state LoginStateStatus) LoginStep {
	return &loginStepImpl{name: name, state: state}
}

type LoginStateStatus string

const (
	Success LoginStateStatus = "success"
	Failed  LoginStateStatus = "failed"
)

type LoginState interface {
	UserID() *string
	ID() *string
	Steps() *[]LoginStep
	AddStep(loginStep LoginStep) LoginState
	IsSuccess() bool
	IsFail() bool
	State() LoginStateStatus
}

type loginStateImpl struct {
	loginSteps []LoginStep
	userId     string
	id         string
}

func (state *loginStateImpl) ID() *string {
	return &state.id
}

func (state *loginStateImpl) UserID() *string {
	return &state.userId
}

func (state *loginStateImpl) IsSuccess() bool {
	return state.State() == Success
}

func (state *loginStateImpl) IsFail() bool {
	return !state.IsSuccess()
}

func (state *loginStateImpl) State() LoginStateStatus {
	for _, item := range state.loginSteps {
		if item.IsFail() {
			return Failed
		}
	}
	return Success
}

func (state *loginStateImpl) Steps() *[]LoginStep {
	return &state.loginSteps
}

func (state *loginStateImpl) AddStep(loginStep LoginStep) LoginState {
	state.loginSteps = append(state.loginSteps, loginStep)
	return state
}

func NewLoginState(userId uuid.UUID) LoginState {
	return &loginStateImpl{
		loginSteps: []LoginStep{},
		userId: userId.String(),
		id: uuid.NewV4().String(),
	}
}
