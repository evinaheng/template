package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
}

// Send message to slack
func Send(message, slackURL string, withMention bool) error {

	// Add mention @channel
	if withMention {
		message = fmt.Sprintf(`<!channel>\n%s`, message)
	}

	form := url.Values{}
	form.Set("payload", fmt.Sprintf(`{"text": "%s"}`, message))
	rb := strings.NewReader(form.Encode())

	req, _ := http.NewRequest("POST", slackURL, rb)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

// SendAttachment message with attachment to slack
func SendAttachment(attachment Attachment, slackURL string, withMention bool) error {

	// Add mention @channel
	if withMention {
		attachment.Pretext = fmt.Sprintf(`<!channel>\n%s`, attachment.Pretext)
	}

	msg, _ := json.Marshal(attachment)

	form := url.Values{}
	form.Set("payload", fmt.Sprintf(`{"attachments": [%s]}`, string(msg)))
	rb := strings.NewReader(form.Encode())

	req, _ := http.NewRequest("POST", slackURL, rb)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
