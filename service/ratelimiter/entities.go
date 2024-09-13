package ratelimiter


type SetUserConfigParams struct {
	UserID    string `json:"userID"`
	RateLimit int    `json:"rateLimit"`
}

type Product struct {
	ID   int
	Name string
}