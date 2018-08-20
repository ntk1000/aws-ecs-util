package main

import (
	"os"
	"testing"
)

const WEBHOOK = "SLACK_WEBHOOK_URL"

func TestPost(t *testing.T) {

	env := os.Getenv(WEBHOOK)
	os.Setenv(WEBHOOK, "")
	s := NewSlack("", "", "")
	err := s.Post()
	if err == nil {
		t.Errorf(ExitMsg, err, "some error")
	}

	os.Setenv(WEBHOOK, env)
	s = NewSlack("", "", "test dayo from aws-ecs-util :sake:")
	err = s.Post()
	if err != nil {
		t.Errorf(ExitMsg, err, nil)
	}
}

func TestJSON(t *testing.T) {
	s := NewSlack("", "", "")
	json, err := s.HandleJSON()
	if err == nil {
		t.Errorf(ExitMsg, json, "blank string")
		t.Errorf(ExitMsg, err, "text empty error")
	}

	s = NewSlack("test", "", "")
	json, err = s.HandleJSON()
	if err == nil {
		t.Errorf(ExitMsg, json, "blank string")
		t.Errorf(ExitMsg, err, "text empty error")
	}

	s = NewSlack("", "test", "")
	json, err = s.HandleJSON()
	if err == nil {
		t.Errorf(ExitMsg, json, "blank string")
		t.Errorf(ExitMsg, err, "text empty error")
	}

	test_text := "test\n <!channel>"
	s = NewSlack("", "", "test")
	json, err = s.HandleJSON()
	if err != nil || json != `{"text":"`+test_text+`"}` {
		t.Errorf(ExitMsg, json, "text json")
		t.Errorf(ExitMsg, err, nil)
	}

	s = NewSlack("test", "", "test")
	json, err = s.HandleJSON()
	if err != nil || json != `{"username":"test","text":"`+test_text+`"}` {
		t.Errorf(ExitMsg, json, "username, text json")
		t.Errorf(ExitMsg, err, nil)
	}

	s = NewSlack("", "test", "test")
	json, err = s.HandleJSON()
	if err != nil || json != `{"channel":"test","text":"`+test_text+`"}` {
		t.Errorf(ExitMsg, json, "channel, text json")
		t.Errorf(ExitMsg, err, nil)
	}

	s = NewSlack("test", "test", "test")
	json, err = s.HandleJSON()
	if err != nil || json != `{"channel":"test","username":"test","text":"`+test_text+`"}` {
		t.Errorf(ExitMsg, json, "channel, username, text json")
		t.Errorf(ExitMsg, err, nil)
	}

}
