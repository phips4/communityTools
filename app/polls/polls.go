package polls

import (
	"strconv"
	"time"
)

type PollOption struct {
	Option    string `json:"option"`
	VoteCount int    `json:"votesCount"`
}

type Vote struct {
	IP          string `bson:"ip" json:"ip"`
	CookieToken string `bson:"cookieToken" json:"cookie_token"`
	Option 		string `json:"option"`
}

// we marshall the struct to json and not vice versa. Since that, fields annotated with "-" are fine.
type Poll struct {
	ID              string        `bson:"_id" json:"id"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	CreatedAt       time.Time     `bson:"createdAt" json:"created_at,omitempty"`
	CookieCheck     bool          `bson:"cookieCheck" json:"cookie_check"`         //TODO:
	MultipleOptions bool          `bson:"multipleOptions" json:"multiple_options"` //TODO:
	Options         []*PollOption `json:"options,omitempty"`
	Votes           []*Vote       `json:"votes,omitempty"`
	VotingStopped   bool          `json:"votingStopped"`
	EditToken       string        `bson:"editToken" json:"-"`
	DeleteAt        time.Time     `json:"delete_at"` //TODO: add TTL mongodb index
}

func NewPoll(id, title, desc, cookie, multiOptions, editToken, deleteDays string, options []string) *Poll {
	pollOption := make([]*PollOption, len(options))

	for i := range pollOption {
		pollOption[i] = &PollOption{options[i], 0}
	}

	cookieCheck, err := strconv.ParseBool(cookie)
	if err != nil {
		return nil
	}

	multipleOptions, err := strconv.ParseBool(multiOptions)
	if err != nil {
		return nil
	}

	n, err := strconv.ParseUint(deleteDays, 10, 11)
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
		nil,
		false,
		editToken,
		time.Now().Add(time.Hour * time.Duration(24) * time.Duration(n)),
	}
}
