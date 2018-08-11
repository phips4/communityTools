package logic

import "testing"

func TestDefaultValidation(t *testing.T) {
	str := "abc123"
	if DefaultValidation(str) == false {
		t.Errorf("'%s' should be valid.", str)
	}

	str = "abc-123"
	if DefaultValidation(str) == false {
		t.Errorf("'%s' should be valid.", str)
	}

	str = ""
	if DefaultValidation(str) == true {
		t.Errorf("'%s' should not be valid, beacuse it is a empty string.", str)
	}

	str = "db.communityTools.find( { $where: function() { return obj.totalVotes > 10; } } );"
	if DefaultValidation(str) == true {
		t.Error("Should not be valid, because the string is too long.")
	}

	t.Log("DefaultValidation tested.")
}

func TestValidateID(t *testing.T) {
	str := "walrus89"
	if ValidateID(str) == false {
		t.Errorf("'%s' should be valid.", str)
	}

	str = "db.c.find({$where: blub});"
	if ValidateID(str) == true {
		t.Error("Should not be valid, because is contains illegal characters.")
	}

	t.Log("ValidateID tested.")
}

func TestValidatePostParams(t *testing.T) {
	if ValidatePostParams("my title", "a description", "false", "true", []string{"apples", "pears"}) == false {
		t.Error("post params. should be valid.")
	}

	if ValidatePostParams("foo", "bar", "notabool", "", []string{"foo", "bar"}) == true {
		t.Error("bool params are not valid, thus this should be not valid")
	}

	if ValidatePostParams("foo", "bar", "true", "false", []string{}) == true {
		t.Error("options array is not valid, thus this should be not valid")
	}

	if ValidatePostParams("foo", "bar", "true", "false", []string{""}) == true {
		t.Error("options array is not valid, thus this should be not valid")
	}

	if ValidatePostParams("foo", "bar", "true", "false", []string{"apple"}) == true {
		t.Error("options array is not valid, thus this should be not valid")
	}

	t.Log("ValidatePostParams tested.")
}
