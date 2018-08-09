package polls

import "time"

type PollOption struct {
	Option string `json:"option"`
	Votes  int    `json:"votes,omitempty"`
}

type Poll struct {
	ID              string       `bson:"_id" json:"id"`
	Title           string       `bson:"title" json:"title"`
	Description     string       `bson:"description" json:"description"`
	CreatedAt       time.Time    `json:"created_at,omitempty"`
	CookieCheck     bool         `json:"cookie_check"`
	MultipleOptions bool         `json:"multiple_options"`
	Options         []PollOption `json:"options,omitempty"`
}
