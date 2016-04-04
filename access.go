package kuruvi

import r "github.com/venkat/ratelimiter"
import "math/rand"
import "time"

//Access - Combination of Auth and the rate limiters per endpoint.
//Complete info. needed for accessing twitter API.
//There can be many Auths
type Access struct {
	authenticator Authenticator
	RateLimiters  map[string]*r.RateLimiter
}

func NewAccess(authenticator Authenticator, rateLimits map[string]*RateLimit) *Access {
	a := &Access{authenticator: authenticator}
	a.setupRateLimiters(rateLimits)
	rand.Seed(time.Now().UTC().UnixNano())
	return a
}

func (a *Access) setupRateLimiters(rateLimits map[string]*RateLimit) {
	a.RateLimiters = make(map[string]*r.RateLimiter, len(rateLimits))
	for endPoint, rateLimit := range rateLimits {
		authType := a.authenticator.GetAuthType()
		a.RateLimiters[endPoint] = r.NewRateLimiter(
			rateLimit.GetQuota(authType),
			rateLimit.GetRate(authType))
	}
}

//TODO: switch to default ratelimiter if none available. Choose it using
//a method in access
func (a *Access) GetRateLimiter(endPoint string) *r.RateLimiter {
	return a.RateLimiters[endPoint]
}
