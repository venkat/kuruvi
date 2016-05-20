package kuruvi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

//RateLimit has per endpoint ratelimit details
type RateLimit struct {
	Method    string
	EndPoint  string
	Category  string
	UserLimit int
	AppLimit  int
}

//GetRateLimits reads twitter API endpoint ratelimits from a json file
// into an array of RateLimits.
//Take a look at rate_limits.json for the file format and data.
func GetRateLimits(f *os.File) (rateLimits []*RateLimit) {

	jsonBlob, _ := ioutil.ReadAll(f)
	err := json.Unmarshal(jsonBlob, &rateLimits)
	if err != nil {
		log.Fatal(err)
	}
	return
}
