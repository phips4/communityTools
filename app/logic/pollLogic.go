package logic

import (
	"regexp"
	"strconv"
	"github.com/phips4/communityTools/app/polls"
	"github.com/phips4/communityTools/app/utils"
	"net/http"
)

const MaxStringLength = 64
const MaxLongStringLentgh = 1024

func ValidateID(id string) bool {

	if !DefaultValidation(id) {
		return false
	}
	ok, _ := regexp.Match("^[a-zA-Z0-9_]*$", []byte(id))
	return ok
}

func ValidatePostParams(title, desc, cookieCheck, multiOptions string, options []string) bool {

	if len(title) > MaxLongStringLentgh || len(title) < 1 {
		return false
	}

	if len(desc) > MaxLongStringLentgh || len(desc) < 1 {
		return false
	}

	if _, err := strconv.ParseBool(cookieCheck); err != nil {
		return false
	}

	if _, err := strconv.ParseBool(multiOptions); err != nil {
		return false
	}

	if options == nil || len(options) < 2 {
		return false
	}

	for _, elem := range options {
		if !DefaultValidation(elem) {
			return false
		}
	}

	return true
}

func DefaultValidation(str string) bool {
	l := len(str)
	if l > MaxStringLength || l < 1 {
		return false
	}
	return true
}

// returns true if it was successfully
func ApplyVote(poll *polls.Poll, ip, cookieToken, option string) bool {
	optionContains := false
	for _, e := range poll.Options {
		if e.Option == option {
			optionContains = true
			break
		}
	}

	if !optionContains {
		return false
	}

	if poll.Votes == nil {
		votes := make([]*polls.Vote, 1)
		votes[0] = &polls.Vote{IP: ip, CookieToken: cookieToken}
		poll.Votes = votes
		poll.TotalVotes++

		for i := range poll.Options {
			if poll.Options[i].Option == option {
				poll.Options[i].VoteCount++
			}
		}
		return true
	}

	if utils.ContainsIpOrToken(poll.Votes, ip, cookieToken) {
		return false
	}

	poll.Votes = append(poll.Votes, &polls.Vote{IP: ip, CookieToken: cookieToken})

	for i := range poll.Options {
		if poll.Options[i].Option == option {
			poll.Options[i].VoteCount++
		}
	}

	poll.TotalVotes++
	return true
}

func GetIp(req *http.Request) string {
	runes := []rune(req.RemoteAddr)
	return string(runes[1:len(runes) - 7])
}