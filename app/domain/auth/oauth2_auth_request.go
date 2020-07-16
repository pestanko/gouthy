package auth

// https://tools.ietf.org/html/rfc6749#section-4.1.1
type OAuth2AuthorizationRequest struct {
	ClientId      string   `json:"client_id"`
	ResponseType  string   `json:"response_type"`
	RedirectUri   string   `json:"redirect_uri"`
	Scopes        []string `json:"scopes"`
	State         string   `json:"state"`
	Nonce         string   `json:"nonce"` // OPENID
	PKCEChallenge string   `json:"pkce_code_challenge"` // PKCE
	PKCEMethod    string   `json:"pkce_code_method"` // PKCE
}



