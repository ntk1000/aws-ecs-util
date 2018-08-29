package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type Slack struct {
	Name        string            `json:"name"`
	Text        string            `json:"text"`
	Channel     string            `json:"channel"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Fallback   string       `json:"fallback"`
	Color      string       `json:"color"`
	PreText    string       `json:"pretext"`
	AuthorName string       `json:"author_name"`
	AuthorLink string       `json:"author_link"`
	AuthorIcon string       `json:"author_icon"`
	Title      string       `json:"title"`
	TitleLink  string       `json:"title_link"`
	Text       string       `json:"text"`
	Fields     []SlackField `json:"fields"`
	ImageUrl   string       `json:"image_url"`
	ThumbUrl   string       `json:"thumb_url"`
	Footer     string       `json:"footer"`
	FooterIcon string       `json:"footer_icon"`
	Ts         int          `json:"ts"`
}

type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

var WebHookURL string

func (s *Slack) Init() error {
	WebHookURL = os.Getenv("SLACK_WEBHOOK_URL")
	if WebHookURL == "" {
		return errors.New("env SLACK_WEBHOOL_URL is empty")
	}
	return nil
}

func (s *Slack) Setup(t string) {
	s.Text = t
}

func (s *Slack) Post() error {

	j, j_err := json.Marshal(s)
	if j_err != nil {
		return j_err
	}
	//os.Stdout.Write(j)

	req, req_err := http.NewRequest(
		"POST",
		WebHookURL,
		bytes.NewBuffer([]byte(j)),
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
