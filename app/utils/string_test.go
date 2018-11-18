package utils

import (
	"testing"
)

func TestHasPrefix(t *testing.T) {

	if HasPrefix("img.community.tools", []string{"abc", "banana", ".img"}) {
		t.Error("expected to be false, no suitable prefix is in the array.")
	}

	if !HasPrefix("paste.community.tools", []string{"orange", ".com", "paste.", "asdf.", "▇▇"}) {
		t.Error("expected to be true, a suitable prefix is in the array.")
	}

	if HasPrefix("", []string{"github.com"}) {
		t.Error("expected to be false, no string to lookup.")
	}

	if HasPrefix("github.com", nil) {
		t.Error("expected to be false, no array given.")
	}

	if HasPrefix("i am a string", []string{""}) {
		t.Error("expected to be false, empty string is not a good prefix.")
	}
}

func TestRandomStringFrom(t *testing.T) {
	rndBinStr := RandomStringFrom(12, []rune{'0', '1'})
	if len(rndBinStr) != 12 {
		t.Error("expected to be true, length should be 12.")
	}

	for _, r := range rndBinStr {
		if r != '0' && r != '1' {
			t.Error("expected to be true, random string ")
		}
	}

	rndStr := RandomString(43)
	if len(rndStr) != 43 {
		t.Error("expected to be true, length should be 43.")
	}

	if RandomStringFrom(-2, []rune{}) != "" {
		t.Error("expected to be false, invalid params given")
	}

	if RandomStringFrom(40, []rune{}) != "" {
		t.Error("expected to be false, invalid params given")
	}
}
