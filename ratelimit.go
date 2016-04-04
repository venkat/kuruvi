package kuruvi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

//RateLimit per endpoint ratelimit details
type RateLimit struct {
	Method    string
	EndPoint  string
	Category  string
	Quota int
	Window    time.Duration
}

//NewRateLimit Create a new ratelimit object
func NewRateLimit(method, endPoint, category string, quota int, window time.Duration) *RateLimit {
	return &RateLimit{method, endPoint, category, quota, window}
}

func (r *RateLimit) GetRate(authType int) time.Duration {
	return r.Window / time.Duration(r.Quota)
}

func GetRateLimits(f *os.File, window time.Duration) map[string]*RateLimit {

	var rateLimits []*RateLimit
	jsonBlob, _ := ioutil.ReadAll(f)
	err := json.Unmarshal(jsonBlob, &rateLimits)
	if err != nil {
		log.Fatal(err)
	}
	r := make(map[string]*RateLimit, len(rateLimits))
	for _, rateLimit := range rateLimits {
		r[rateLimit.EndPoint] = rateLimit
		r[rateLimit.EndPoint].Window = window
	}
	return r
}
