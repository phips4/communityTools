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
}

//
// we marshall the struct to json and not vice versa. Since that, fields annotated with "-" are fine.
type Poll struct {
	ID              string        `bson:"_id" json:"id"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	CreatedAt       time.Time     `bson:"createdAt" json:"created_at,omitempty"`
	CookieCheck     bool          `bson:"cookieCheck" json:"cookie_check"`
	MultipleOptions bool          `bson:"multipleOptions" json:"multiple_options"`
	Options         []*PollOption `json:"options,omitempty"`
	Votes           []*Vote       `json:"votes,omitempty"`
	TotalVotes      int           `json:"totalVotes"`
	VotingStopped   bool          `json:"votingStopped"`
	EditToken       string        `bson:"editToken" json:"-"`
	DeleteAt        time.Time     `json:"delete_at"`
}

func NewPoll(id, title, desc, cookie, multiOptions, editToken string, deleteAt int, options []string) *Poll {
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

	return &Poll{
		id,
		title,
		desc,
		time.Now(),
		cookieCheck,
		multipleOptions,
		pollOption,
		nil,
		0,
		false,
		editToken,
		time.Now().Add(time.Hour * time.Duration(24) * time.Duration(deleteAt)),
	}
}
