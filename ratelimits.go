package kuruvi

//RateLimit has per endpoint ratelimit details
type RateLimit struct {
	Method    string
	EndPoint  string
	Category  string
	UserLimit int
	AppLimit  int
}

var rateLimits = []*RateLimit{
	{Method: "GET", EndPoint: "application/rate_limit_status", Category: "application", UserLimit: 180, AppLimit: 180},
	{Method: "GET", EndPoint: "favorites/list", Category: "favorites", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "followers/ids", Category: "followers", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "followers/list", Category: "followers", UserLimit: 15, AppLimit: 30},
	{Method: "GET", EndPoint: "friends/ids", Category: "friends", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "friends/list", Category: "friends", UserLimit: 15, AppLimit: 30},
	{Method: "GET", EndPoint: "friendships/show", Category: "friendships", UserLimit: 180, AppLimit: 15},
	{Method: "GET", EndPoint: "help/configuration", Category: "help", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "help/languages", Category: "help", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "help/privacy", Category: "help", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "help/tos", Category: "help", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "lists/list", Category: "lists", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "lists/members", Category: "lists", UserLimit: 180, AppLimit: 15},
	{Method: "GET", EndPoint: "lists/members/show", Category: "lists", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "lists/memberships", Category: "lists", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "lists/ownerships", Category: "lists", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "lists/show", Category: "lists", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "lists/statuses", Category: "lists", UserLimit: 180, AppLimit: 180},
	{Method: "GET", EndPoint: "lists/subscribers", Category: "lists", UserLimit: 180, AppLimit: 15},
	{Method: "GET", EndPoint: "lists/subscribers/show", Category: "lists", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "lists/subscriptions", Category: "lists", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "search/tweets", Category: "search", UserLimit: 180, AppLimit: 450},
	{Method: "GET", EndPoint: "statuses/lookup", Category: "statuses", UserLimit: 180, AppLimit: 60},
	{Method: "GET", EndPoint: "statuses/oembed", Category: "statuses", UserLimit: 180, AppLimit: 180},
	{Method: "GET", EndPoint: "statuses/retweeters/ids", Category: "statuses", UserLimit: 15, AppLimit: 60},
	{Method: "GET", EndPoint: "statuses/retweets/:id", Category: "statuses", UserLimit: 15, AppLimit: 60},
	{Method: "GET", EndPoint: "statuses/show/:id", Category: "statuses", UserLimit: 180, AppLimit: 180},
	{Method: "GET", EndPoint: "statuses/user_timeline", Category: "statuses", UserLimit: 180, AppLimit: 300},
	{Method: "GET", EndPoint: "trends/available", Category: "trends", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "trends/closest", Category: "trends", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "trends/place", Category: "trends", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "users/lookup", Category: "users", UserLimit: 180, AppLimit: 60},
	{Method: "GET", EndPoint: "users/show", Category: "users", UserLimit: 180, AppLimit: 180},
	{Method: "GET", EndPoint: "users/suggestions", Category: "users", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "users/suggestions/:slug", Category: "users", UserLimit: 15, AppLimit: 15},
	{Method: "GET", EndPoint: "users/suggestions/:slug/members", Category: "users", UserLimit: 15, AppLimit: 15},
}

//GetRateLimits reads twitter API endpoint ratelimits from a json file
// into an array of RateLimits.
//Take a look at rate_limits.json for the file format and data.
func GetRateLimits() []*RateLimit {
	return rateLimits
}
