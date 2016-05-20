package kuruvi

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type appAuth struct {
	*Auth
	bearerToken string
}

func newAppAuth(key, secret string) *appAuth {
	return &appAuth{Auth: &Auth{ConsumerKey: key, ConsumerSecret: secret}}
}

func (a *appAuth) getBearerToken() {
	client := &http.Client{}
	//Step 1: Encode consumer key and secret
	encodedKeySecret := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
		url.QueryEscape(a.ConsumerKey),
		url.QueryEscape(a.ConsumerSecret))))

	//Step 2: Obtain a bearer token
	//The body of the request must be grant_type=client_credentials
	reqBody := bytes.NewBuffer([]byte(`grant_type=client_credentials`))
	//The request must be a HTTP POST request
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", reqBody)
	if err != nil {
		log.Fatal(err, client, req)
	}
	//The request must include an Authorization header formatted as
	//Basic <base64 encoded value from step 1>.
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedKeySecret))
	//The request must include a Content-Type header with
	//the value of application/x-www-form-urlencoded;charset=UTF-8.
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("Content-Length", strconv.Itoa(reqBody.Len()))

	//Issue the request and get the bearer token from the JSON you get back
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err, resp)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, respBody)
	}

	type BearerToken struct {
		AccessToken string `json:"access_token"`
	}
	var b BearerToken
	json.Unmarshal(respBody, &b)
	a.bearerToken = b.AccessToken
}

func (a *appAuth) setAuthHeader(r *http.Request) {
	if a.bearerToken == "" {
		a.getBearerToken()
	}

	//Step 3: Authenticate API requests with the bearer token
	//include an Authorization header formatted as
	//Bearer <bearer token value from step 2>
	r.Header.Add("Authorization",
		fmt.Sprintf("Bearer %s", a.bearerToken))
}

func (a *appAuth) getAuthType() int {
	return App
}
