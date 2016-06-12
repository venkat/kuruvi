//Package kuruvi is a simple wrapper around Twitter's API for GET enpoints.
//
//It takes care of two things:
//
//1. authentication - it supports Twitter's user auth and app auth
//
//2. rate limits - supports the fine-grained per endpoint rate-limiting specified
//by Twitter at https://dev.twitter.com/rest/public/rate-limits. The throttling
//ensures that you won't get ratelimited by twitter
//
//SetupKuruvi is the entry point. Create a Kuruvi object with that and call Get
//on the object to access Twitter's API.
package kuruvi

import (
	"encoding/json"
	"os"
	"strconv"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	r "github.com/venkat/ratelimiter"
)

const (
	debug = false
)

const (
	//App signifies the app auth type
	App = iota
	//User signifies the user auth type
	User
)

const (
	//UseAppAuth will only use App Auth (and related rate limiting)
	UseAppAuth = iota
	//UseUserAuth will only use User Auth (and related rate limiting)
	UseUserAuth
	//UseBoth will mix App Auth and User Auth (combining the rate limit quota for App and User auth)
	UseBoth
)

//APIBase is the fixed prefix added to all twitter API endpoints
const APIBase = "https://api.twitter.com/1.1/"

//Authenticator helps set authentication header for the given request
//authentication header setting depends on the auth type (app auth or user auth)
type authenticator interface {
	setAuthHeader(r *http.Request)
	getAuthType() int
}

//Auth is a collection of possible app auth and user auth tokens you can get
//from Twitter's developer console. Only set the keys you need for
//the auth type that is being used
type Auth struct {
	ConsumerKey       string `json:"consumerKey"`
	ConsumerSecret    string `json:"consumerSecret"`
	AccessTokenKey    string `json:"accessTokenKey"`
	AccessTokenSecret string `json:"accessTokenSecret"`
}

//GetAuthKeys reads authenticaton keys from a json file into an Auth struct
//check out template_auth.json for a sample of how the file must look like
//only set the Auth tokens you have for the auth type (app or user auth) you are using
func GetAuthKeys(f *os.File) (auth *Auth) {
	jsonBlob, _ := ioutil.ReadAll(f)
	err := json.Unmarshal(jsonBlob, &auth)
	if err != nil {
		log.Fatal(err)
	}
	return
}

//Access has an Authenticator and the set of per endpoint ratelimiters, all the information
//needed to access Twitter's APIs without getting ratelimited
type access struct {
	Auth         authenticator
	rateLimiters map[string]*r.RateLimiter
}

//NewAccess creates a new access object with the authentication information and the ratelimit
//information. window is the time window used by twitter for rate limiting (currently 15 mins).
func newAccess(auth authenticator, rateLimits []*RateLimit, window time.Duration) *access {
	a := &access{Auth: auth}
	a.setupRateLimiters(rateLimits, window)
	return a
}

func (a *access) accessTypeName() string {
	if a.Auth.getAuthType() == App {
		return "app"
	}

	return "user"
}
func (a *access) setupRateLimiters(rateLimits []*RateLimit, window time.Duration) {
	a.rateLimiters = make(map[string]*r.RateLimiter, len(rateLimits))
	authType := a.Auth.getAuthType()
	for _, rateLimit := range rateLimits {

		quota := rateLimit.AppLimit
		if authType == User {
			quota = rateLimit.UserLimit
		}

		//quota-- //making one less API call because twitter doesn't roll over its timewindow correctly by 15min boundary

		a.rateLimiters[rateLimit.EndPoint] = r.NewRateLimiter(
			quota,
			window/time.Duration(quota),
			fmt.Sprintf("%s - %s", a.accessTypeName(), rateLimit.EndPoint))
	}
}

func (a *access) getRateLimiter(endPoint string) *r.RateLimiter {
	var endPointTemplate string

	switch {
	case strings.HasPrefix(endPoint, "statuses/retweets/") == true:
		endPointTemplate = "statuses/retweets/:id"
	case strings.HasPrefix(endPoint, "statuses/show/") == true:
		endPointTemplate = "statuses/show/:id"
	case strings.HasPrefix(endPoint, "users/suggestions/") == true && strings.HasSuffix(endPoint, "members") == true:
		endPointTemplate = "users/suggestions/:slug/members"
	case strings.HasPrefix(endPoint, "users/suggestions/") == true:
		endPointTemplate = "users/suggestions/:slug"
	default:
		endPointTemplate = endPoint
	}

	return a.rateLimiters[endPointTemplate]
}

//Kuruvi works with either an AppAccess or an UserAccess or both
type Kuruvi struct {
	AppAccess  *access
	UserAccess *access
}

//SetupKuruvi sets up the API client with the ratelimit window, auth information
//The last parameter specifies what kind of authentication to use
//It should be one of:
//UseAppAuth - only use App Auth (and related rate limiting)
//UseUserAuth - only use User Auth (and related rate limiting)
//UseBoth - mix App Auth and User Auth (combining the rate limit quota for App and User auth)
func SetupKuruvi(window time.Duration, authKeys *Auth, useAuthType int) *Kuruvi {
	var appAccess *access
	var userAccess *access

	switch useAuthType {
	case UseBoth:
		{
			appAccess = newAccess(newAppAuth(authKeys.ConsumerKey, authKeys.ConsumerSecret), rateLimits, window)

			userAccess = newAccess(newUserAuth(authKeys.ConsumerKey,
				authKeys.ConsumerSecret,
				authKeys.AccessTokenKey,
				authKeys.AccessTokenSecret), rateLimits, window)

		}
	case UseAppAuth:
		appAccess = newAccess(newAppAuth(authKeys.ConsumerKey, authKeys.ConsumerSecret), rateLimits, window)
	case UseUserAuth:
		userAccess = newAccess(newUserAuth(authKeys.ConsumerKey,
			authKeys.ConsumerSecret,
			authKeys.AccessTokenKey,
			authKeys.AccessTokenSecret), rateLimits, window)
	}

	//You can just have one of these two Accesses. The other can be nil.
	return newKuruvi(appAccess, userAccess)

}

//NewKuruvi gives a Kuruvi object which is the API wrapper
//You can pass in both appAccess (app auth) and userAccess (user auth)
//you can pass in only one of them but at least one of them is needed
func newKuruvi(appAccess *access, userAccess *access) *Kuruvi {
	k := &Kuruvi{AppAccess: appAccess, UserAccess: userAccess}
	return k
}

func (k *Kuruvi) prepare(req *http.Request, endPoint string) {
	switch {
	case k.AppAccess == nil && k.UserAccess == nil:
		log.Fatal("must have atleeast user access or app access for Kuruvi to work")
	case k.AppAccess == nil:
		k.UserAccess.Auth.setAuthHeader(req)
		k.UserAccess.getRateLimiter(endPoint).Throttle()
	case k.UserAccess == nil:
		k.AppAccess.Auth.setAuthHeader(req)
		k.AppAccess.getRateLimiter(endPoint).Throttle()
	default:
		appTokensLeft := k.AppAccess.getRateLimiter(endPoint).TokensLeft()
		userTokensLeft := k.UserAccess.getRateLimiter(endPoint).TokensLeft()
		if appTokensLeft > userTokensLeft {
			k.AppAccess.Auth.setAuthHeader(req)
			k.AppAccess.getRateLimiter(endPoint).Throttle()
		} else if userTokensLeft > appTokensLeft {
			k.UserAccess.Auth.setAuthHeader(req)
			k.UserAccess.getRateLimiter(endPoint).Throttle()
		} else {
			k.UserAccess.Auth.setAuthHeader(req)
			k.UserAccess.getRateLimiter(endPoint).Throttle()
		}
	}
}

//Get makes a GET for the given API endPoint, form has the GET parameters
//to get the JSON response in bytes
func (k *Kuruvi) Get(endPoint string, form url.Values) ([]byte, error) {
	log.SetFlags(log.Lshortfile)

	req, err := http.NewRequest("GET", APIBase+endPoint+".json", nil)
	if err != nil {
		log.Fatal(err, req)
	}

	req.URL.RawQuery = form.Encode()

	client := &http.Client{}
	k.prepare(req, endPoint)
	log.Printf("calling Get for %v and time is %v", endPoint, time.Now())
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err, resp)
	}

	resetTimestamp, _ := strconv.Atoi(resp.Header.Get("X-Rate-Limit-Reset"))
	log.Printf("response headers for Get for %v at %v. Rate Limit %v. Rate Limit Remaining %v. Rate Limit Resets at %v. Time Remaining %v.",
		endPoint,
		time.Now(),
		resp.Header.Get("X-Rate-Limit-Limit"),
		resp.Header.Get("X-Rate-Limit-Remaining"),
		time.Unix(int64(resetTimestamp), 0),
		time.Unix(int64(resetTimestamp), 0).Sub(time.Now()))

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, body)
	}

	if debug {
		out, err := httputil.DumpRequestOut(req, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(strings.Replace(string(out), "\r", "", -1))
		fmt.Println(string(body))
		out, err = httputil.DumpResponse(resp, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(strings.Replace(string(out), "\r", "", -1))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP GET unsuccessful. endpoint: %s values: %s respCode: %d", endPoint, form, resp.StatusCode)
	}

	return body, nil
}
