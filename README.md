Kuruvi is a go package that is a simple wrapper for the GET endpoints of Twitter's API.

Kuruvi takes care of two things:

1. **authentication** - it supports Twitter's [user auth](https://dev.twitter.com/oauth) and [app auth](https://dev.twitter.com/oauth/application-only) 
2. **rate limits** - supports the [fine-grained per endpoint rate-limiting](https://dev.twitter.com/rest/public/rate-limits) specified by Twitter. The throttling ensures that you won't get ratelimited by Twitter.

Unmarshaling the JSON response is *not* in the scope for this wrapper. If you are looking for a full-fledged API client, I recommend [anaconda](https://github.com/ChimeraCoder/anaconda) by [ChimeraCoder](https://github.com/ChimeraCoder/).

Kurvui currently does not support streaming endpoints.

Using Kuruvi, I was able to hit different twitter enpoints for 8 hours straight without getting ratelimited.

## Using Kuruvi

Here is an example of how to setup the API client and GET a Twitter user object.

```go
package main

import (
        "encoding/json"
        "fmt"
        "log"
        "net/url"
        "os"
        "time"

        "github.com/venkat/kuruvi"
)

//User is used when unmarshalling Twitter's JSON response
type User struct {
        ID          int64  `json:"id"`
        Name        string `json:"name"`
        ScreenName  string `json:"screen_name"`
        Description string `json:"description"`
        Protected   bool   `json:"protected"`
}

//helper function to get an open file handle
func getFile(fileName string) *os.File {
        f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
        if err != nil {
                log.Fatal(err)
        }
        return f
}

func main() {

        //auth.json contains the authentication tokens you get from Twitter's
        //Application Management portal - apps.twitter.com
        //Checkout template_auth.json for the format
        authFile := getFile("auth.json")

        //Twitter has a 15 minute rate limit rollover window
        //adding an extra minute to avoid edge cases with twitter rolling over its time window
        window := 15*time.Minute + 1*time.Minute

        //Setup the API client with the ratelimit window and auth information
        //The last parameter specifies what kind of authentication to use
        //It should be one of:
        //UseAppAuth - only use App Auth (and related rate limiting)
        //UseUserAuth - only use User Auth (and related rate limiting)
        //UseBoth - mix App Auth and User Auth (combining the rate limit quota for App and User auth)
        k := kuruvi.SetupKuruvi(
                window,
                kuruvi.GetAuthKeys(authFile),
                kuruvi.UseBoth)

        //check out Twitter's API documentation for all available enpoints
        //and their parameters at: https://dev.twitter.com/rest/public
        v := url.Values{}
        v.Add("screen_name", "annacoder")
        data, err := k.Get("users/show", v)
        if err != nil {
                log.Fatal(err)
        }

        var u *User
        err = json.Unmarshal(data, &u)
        if err != nil {
                log.Fatal(err)
        }
        fmt.Println(u)
}
```

## TODO

1. Write tests.
2. Support POST endpoints.
3. Add more examples.
4. Streaming endpoints support.
