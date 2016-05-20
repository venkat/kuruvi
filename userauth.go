package kuruvi

import (
	"log"
	"net/http"

	"github.com/garyburd/go-oauth/oauth"
)

type userAuth struct {
	*Auth
	oauthClient oauth.Client
}

func newUserAuth(consumerKey, consumerSecret, accessTokenKey, accessTokenSecret string) *userAuth {
	u := &userAuth{Auth: &Auth{
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

func (u *userAuth) setAuthHeader(req *http.Request) {
	u.oauthClient.Credentials.Token = u.ConsumerKey
	u.oauthClient.Credentials.Secret = u.ConsumerSecret
	credentials := &oauth.Credentials{Token: u.AccessTokenKey, Secret: u.AccessTokenSecret}
	err := u.oauthClient.SetAuthorizationHeader(req.Header, credentials, req.Method, req.URL, req.Form)
	if err != nil {
		log.Fatal(err, req)
	}
}

func (u *userAuth) getAuthType() int {
	return User
}
