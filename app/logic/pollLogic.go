package logic

import (
	"errors"
	"github.com/phips4/communityTools/app/entity"
	"github.com/phips4/communityTools/app/utils"
	"strconv"
)

var (
	ErrAlreadyVoted   = errors.New("you have already voted")
	ErrSomethingWrong = errors.New("something went wrong")
)

func ValidatePostParams(title, desc, cookieCheck, multiOptions, deleteDays string, options []string) bool {

	if len(title) > MaxLongStringLength || len(title) < 1 {
		return false
	}

	if len(desc) > MaxLongStringLength || len(desc) < 1 {
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
		if !IdLengthValidation(elem) {
			return false
		}
	}

	// we only allow ints in range of 0 to 2047.
	// So we don't waste all the other 53 bytes and 2047 days are enough I think. That are 5.6 years
	n, err := strconv.ParseUint(deleteDays, 10, 11)
	if err != nil || n < 1 {
		return false
	}

	return true
}

// returns nil if it was successfully
func ApplyVote(poll *entity.Poll, ip, cookieToken, option string) error {
	optionContains := false
	for _, e := range poll.Options {
		if e.Option == option {
			optionContains = true
			break
		}
	}

	if !optionContains {
		return ErrSomethingWrong
	}

	apply := func() {
		for i := range poll.Options {
			if poll.Options[i].Option == option {
				poll.Options[i].VoteCount++
				poll.Options[i].Option = option
			}
		}
	}

	if poll.Votes == nil {
		votes := make([]*entity.Vote, 1)
		votes[0] = &entity.Vote{IP: ip, CookieToken: cookieToken, Option: option}
		poll.Votes = votes
		apply()
		return nil
	}

	if poll.MultipleOptions {
		if utils.VotedFor(option, poll.Votes, ip, cookieToken) {
			return ErrAlreadyVoted
		} else {
			apply()
		}
	} else if utils.ContainsIpOrToken(poll.Votes, ip, cookieToken) {
		return ErrAlreadyVoted
	}

	poll.Votes = append(poll.Votes, &entity.Vote{IP: ip, CookieToken: cookieToken, Option: option})
	apply()

	return nil
}