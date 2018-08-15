package utils

import (
	"testing"
	"github.com/phips4/communityTools/app/polls"
)

func TestContains(t *testing.T) {
	src := []string{"apple", "banana", "avocado"}

	if Contains(src, "apple") == false {
		t.Error("apple could not be found but it exists.")
	}

	src = []string{"avocado"}
	if Contains(src, "apple") == true {
		t.Error("apple could be found but it does not exist.")
	}

	src = []string{}

	if Contains(src, "apple") == true {
		t.Error("apple could be found but it does not exist")
	}

}

func TestContainsMany(t *testing.T) {
	src := []string{"apple", "coconut", "feijoa", "kiwi", "pear"}
	search := []string{"feijoa", "kiwi", "pear"}

	if ContainsMany(src, search) == false {
		t.Errorf("elements: %v could not be found in array: %v", search, src)
	}
}

func TestContainsIpOrToken(t *testing.T) {
	votes := make([]*polls.Vote, 3)
	// random generated IPs and tokens
	votes[0] = &polls.Vote{IP: "174.26.249.79", CookieToken: "POG74ViLCAAHDOXycdnE6CmkY"}
	votes[1] = &polls.Vote{IP: "221.200.87.228", CookieToken: "sSKrakmHNPbap0F31GFWT"}
	votes[2] = &polls.Vote{IP: "181.184.85.10", CookieToken: "eY38Ut58if1eqGE8lVKX"}

	if ContainsIpOrToken(votes, "174.26.249.79", "no token") == false {
		t.Errorf("slice contains ip: '%s' but it fails.", "174.26.249.79")
	}

	if ContainsIpOrToken(votes, "210.137.16.208", "eY38Ut58if1eqGE8lVKX") == false {
		t.Errorf("slice contains token: '%s' but it fails.", "eY38Ut58if1eqGE8lVKX")
	}

	if ContainsIpOrToken(votes, "174.26.249.79", "POG74ViLCAAHDOXycdnE6CmkY") == false {
		t.Errorf("slice contains ip and token: '%s' and '%s' but it fails.", "174.26.249.79", "POG74ViLCAAHDOXycdnE6CmkY")
	}

	if ContainsIpOrToken(votes, "", "") == true {
		t.Errorf("slice does not contain any empty string but it does not fail.")
	}
}