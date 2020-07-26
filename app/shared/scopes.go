package shared

import "github.com/pestanko/gouthy/app/shared/utils"

const (
	ScopeSession = "session"
	ScopeUI = "ui"
	ScopeUnauthorized = "unauthorized"
	ScopeAnonymous = "anon"
)

type Scopes []string


func (s Scopes) IsSession() bool {
	return utils.StringSliceContains(s, ScopeSession)
}


func (s Scopes) IsUI() bool {
	return utils.StringSliceContains(s, ScopeUI)
}

func (s Scopes) IsUnauthorized() bool {
	return utils.StringSliceContains(s, ScopeUnauthorized)
}

func (s Scopes) IsAnonymous() bool {
	return utils.StringSliceContains(s, ScopeUnauthorized)
}





