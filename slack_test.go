package main

import (
	"os"
	"testing"
)

const WEBHOOK = "SLACK_WEBHOOK_URL"

func TestPost(t *testing.T) {

	env := os.Getenv(WEBHOOK)
	os.Setenv(WEBHOOK, "")
	s := &Slack{}
	err := s.Init()
	if err == nil {
		t.Errorf(ExitMsg, err, "some error")
	}

	os.Setenv(WEBHOOK, env)
	s = &Slack{}
	err = s.Init()
	if err != nil {
		t.Errorf(ExitMsg, err, nil)
	}
	s.Setup("test dayo from aws-ecs-util :sake:")
	err = s.Post()
	if err != nil {
		t.Errorf(ExitMsg, err, nil)
	}
}
