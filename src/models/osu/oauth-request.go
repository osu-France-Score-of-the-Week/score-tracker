package osu

import "net/url"

type OAuthRequest struct {
	ClientID     string
	ClientSecret string
	GrantType    string
	Scope        string
}

func (r OAuthRequest) ToValues() url.Values {
	form := url.Values{}
	form.Set("client_id", r.ClientID)
	form.Set("client_secret", r.ClientSecret)
	form.Set("grant_type", r.GrantType)
	form.Set("scope", r.Scope)
	return form
}
