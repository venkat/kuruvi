package kuruvi

import (
	"github.com/garyburd/go-oauth/oauth"
	"log"
	"net/http"
)

type UserAuth struct {
	*Auth
	oauthClient oauth.Client
}

func NewUserAuth(consumerKey, consumerSecret, accessTokenKey, accessTokenSecret string) *UserAuth {
	u := &UserAuth{Auth: &Auth{
		ConsumerKey:       consumerKey,
		ConsumerSecret:    consumerSecret,
		AccessTokenKey:    accessTokenKey,
		AccessTokenSecret: accessTokenSecret,
	}}

	u.oauthClient = oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
	}

	return u
}

func (u *UserAuth) SetAuthHeader(req *http.Request) {
	u.oauthClient.Credentials.Token = u.ConsumerKey
	u.oauthClient.Credentials.Secret = u.ConsumerSecret
	credentials := &oauth.Credentials{Token: u.AccessTokenKey, Secret: u.AccessTokenSecret}
	err := u.oauthClient.SetAuthorizationHeader(req.Header, credentials, req.Method, req.URL, req.Form)
	if err != nil {
		log.Fatal(err, req)
	}
}

func (u *UserAuth) GetAuthType() int {
	return User
}
