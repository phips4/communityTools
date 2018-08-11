package polls

import (
	"strconv"
	"time"
)

type PollOption struct {
	Option string `json:"option"`
	Votes  int    `json:"votes,omitempty"`
}

type Poll struct {
	ID              string       `bson:"_id" json:"id"`
	Title           string       `bson:"title" json:"title"`
	Description     string       `bson:"description" json:"description"`
	CreatedAt       time.Time    `bson:"createdAt" json:"created_at,omitempty"`
	CookieCheck     bool         `bson:"cookieCheck" json:"cookie_check"`
	MultipleOptions bool         `bson:"multipleOptions" json:"multiple_options"`
	Options         []PollOption `json:"options,omitempty"`
}

func NewPoll(id, title, desc string, cookie, multiOptions string, options []string) *Poll {
	pollOption := make([]PollOption, len(options))

	for i := range pollOption {
		// don't use the element from the range iteration, it's a copy
		pollOption[i].Option = options[i]
		pollOption[i].Votes = 0 //not needed, I know
	}

	cookieCheck, err := strconv.ParseBool(cookie)
	if err != nil {
		return nil
	}

	multipleOptions, err := strconv.ParseBool(multiOptions)
	if err != nil {
		return nil
	}

	return &Poll{
		id,
		title,
		desc,
		time.Now(),
		cookieCheck,
		multipleOptions,
		pollOption,
	}
}
