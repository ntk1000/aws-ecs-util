package main

import (
	"bytes"
	"errors"
	"net/http"
	"os"
)

type Slack struct {
	Name, Text, Channel, WebHookURL string
}

func NewSlack(name, channel, text string) *Slack {
	return &Slack{
		Name:       name,
		Channel:    channel,
		Text:       text,
		WebHookURL: os.Getenv("SLACK_WEBHOOK_URL"),
	}
}

func (s *Slack) Post() error {

	if s.WebHookURL == "" {
		return errors.New("env SLACK_WEBHOOL_URL is empty")
	}

	json, json_err := s.HandleJSON()
	if json_err != nil {
		return json_err
	}

	req, req_err := http.NewRequest(
		"POST",
		s.WebHookURL,
		bytes.NewBuffer([]byte(json)),
	)

	if req_err != nil {
		return req_err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, resp_err := client.Do(req)

	if resp_err != nil {
		return resp_err
	}

	defer resp.Body.Close()
	return nil
}

func (s *Slack) HandleJSON() (string, error) {

	if s.Text == "" {
		return "", errors.New("text is empty")
	} else {
		s.Text += "\n <!channel>"
	}

	if s.Channel == "" && s.Name == "" {
		return `{"text":"` + s.Text + `"}`, nil
	} else if s.Channel == "" && s.Name != "" {
		return `{"username":"` + s.Name + `","text":"` + s.Text + `"}`, nil
	} else if s.Channel != "" && s.Name == "" {
		return `{"channel":"` + s.Channel + `","text":"` + s.Text + `"}`, nil
	} else {
		return `{"channel":"` + s.Channel + `","username":"` + s.Name + `","text":"` + s.Text + `"}`, nil
	}

}
