package logic

import (
	"github.com/phips4/communityTools/app/polls"
	"testing"
	"time"
)

func TestDefaultValidation(t *testing.T) {
	t.Parallel()

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
}

func TestValidateID(t *testing.T) {
	t.Parallel()

	str := "walrus89"
	if ValidateID(str) == false {
		t.Errorf("'%s' should be valid.", str)
	}

	str = "db.c.find({$where: blub});"
	if ValidateID(str) == true {
		t.Error("Should not be valid, because is contains illegal characters.")
	}
}

func TestValidatePostParams(t *testing.T) {
	t.Parallel()

	if ValidatePostParams("my title", "a description", "false", "true", "7", []string{"apples", "pears"}) == false {
		t.Error("post params. should be valid.")
	}

	if ValidatePostParams("foo", "bar", "notabool", "", "7", []string{"foo", "bar"}) == true {
		t.Error("bool params are not valid, thus this should be not valid")
	}

	if ValidatePostParams("foo", "bar", "true", "false", "7", []string{}) == true {
		t.Error("options array is not valid, thus this should be not valid")
	}

	if ValidatePostParams("foo", "bar", "true", "false", "7", []string{""}) == true {
		t.Error("options array is not valid, thus this should be not valid")
	}

	if ValidatePostParams("foo", "bar", "true", "false", "7", []string{"apple"}) == true {
		t.Error("options array is not valid, thus this should be not valid")
	}
}

func TestApplyVote(t *testing.T) {
	t.Parallel()

	p := polls.NewPoll("pollId", "pollTitle", "pollDesc", "true", "true", "editToken", "7", []string{"vote1", "vote2"})
	p.CreatedAt = time.Date(2018, 8, 15, 0, 0, 0, 0, time.UTC)

	if err := ApplyVote(p, "127.0.0.1", "cookieToken", "vote1"); err != nil {
		t.Error("Vote with ip '127.0.0.1' and cookieToken 'cookieToken' cant be applied to struct for vote 'vote1', but it should")
	}

	if err := ApplyVote(p, "192.168.1.1", "verySecretToken", "vote2"); err != nil {
		t.Error("Vote with ip '192.168.1.1' and cookieToken 'verySecretToken' cant be applied to struct for vote 'vote2', but it should be.")
	}

	if err := ApplyVote(p, "192.168.1.1", "verySecretTokenOtherPc", "vote2"); err == nil {
		t.Error("Vote with ip '192.168.1.1' and cookieToken 'verySecretTokenOtherPc' can be applied to struct for vote 'vote2'.")
		t.Error("The IP has already voted, so this should not be possible.")
	}

	if err := ApplyVote(p, "192.168.1.3", "walrus123", "vote2"); err != nil {
		t.Error("Vote with ip '192.168.1.3' and cookieToken 'walrus123' cant be applied to struct for vote 'vote2', but it should")
	}

	if err := ApplyVote(p, "192.168.1.4", "walrus123", "vote1"); err == nil {
		t.Error("Vote with ip '192.168.1.3' and cookieToken 'walrus123' can be applied to struct for vote 'vote1'.")
		t.Error("The cookieToken has already voted, so this should not be possible.")
	}

	if p.Votes == nil {
		t.Error("Votes can not be empty but it is.")
	}

	/*
		out, err := json.MarshalIndent(p, "", "  ")
		if err != nil {
			t.Errorf("err is not nil: %v", err.Error())
		}*/
}
