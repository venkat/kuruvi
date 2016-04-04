//kuruvi package is a minimal API client for twitter1.1 which takes care of
//Twitter's per end-point rate limits and can be used with non-interative
//authentication contexts, namely developer auth or application auth. You
//can feed it multiple auth tokens and it will combine all of them to get
//a higher API rate

package kuruvi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
)

const (
	Debug = false
)

const (
	App = iota
	User
)

const APIBase = "https://api.twitter.com/1.1/"

type Authenticator interface {
	SetAuthHeader(r *http.Request)
	GetAuthType() int
}

type Auth struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessTokenKey    string
	AccessTokenSecret string
}

type Kuruvi struct {
	Accesses []*Access
}

func NewKuruvi(accesses []*Access) *Kuruvi {
	return &Kuruvi{accesses}
}

func (k *Kuruvi) pickAccess(endPoint string) *Access {
	//the one with smaller quota will be depleted faster. Pick weighted by their quota?
	return k.Accesses[rand.Intn(len(k.Accesses))]
}

func (k *Kuruvi) ApplyBestAccess(req *http.Request, endPoint string) {
	cases := make([]reflect.SelectCase, len(k.Accesses))
	for i, a := range k.Accesses {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(a.GetRateLimiter(endPoint).GetThrottleChannel()),
		}
	}

	chosen, _, _ := reflect.Select(cases)
	access := k.Accesses[chosen]
	access.authenticator.SetAuthHeader(req)
}

func (k *Kuruvi) Get(endPoint string, form url.Values, data interface{}) error {
	log.SetFlags(log.Lshortfile)
	req, err := http.NewRequest("GET", APIBase+endPoint+".json", nil)
	if err != nil {
		log.Fatal(err, req)
	}
	req.URL.RawQuery = form.Encode()
	k.ApplyBestAccess(req, endPoint)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err, resp)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, body)
	}

	if Debug {
		out, err := httputil.DumpRequestOut(req, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(strings.Replace(string(out), "\r", "", -1))
		fmt.Println(string(body), "\n")
		out, err = httputil.DumpResponse(resp, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(strings.Replace(string(out), "\r", "", -1))
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("HTTP GET unsuccessful. endpoint: %s values: %s respCode: %s", endPoint, form, resp.StatusCode))
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		log.Fatal(err, string(body))
	}

	return nil
}
