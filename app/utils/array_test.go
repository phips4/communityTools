package utils

import "testing"

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

	//TODO:


}