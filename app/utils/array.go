package utils

import (
	"github.com/phips4/communityTools/app/polls"
)

func ContainsMany(src, content []string) bool {
	for _, contentElem := range content {
		if !Contains(src, contentElem) {
			return false
		}
	}

	return true
}

func Contains(src []string, elem string) bool {

	for _, e := range src {
		if e == elem {
			return true
		}
	}

	return false
}

func ContainsIpOrToken(votes []*polls.Vote, ip, token string) bool {
	for _, v := range votes {
		if v.IP == ip || v.CookieToken == token {
			return true
		}
	}

	return false
}

func VotedFor(option string, votes []*polls.Vote, ip, token string) bool {
	for _, v := range votes {
		if (v.IP == ip || v.CookieToken == token) && v.Option == option {
			return true
		}
	}
	return false
}
