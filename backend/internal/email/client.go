package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const resendEndpoint = "https://api.resend.com/emails"

type Client struct {
	apiKey string
	from   string
	http   *http.Client
}

func NewClient(apiKey, from string) *Client {
	return &Client{
		apiKey: apiKey,
		from:   from,
		http:   &http.Client{Timeout: 10 * time.Second},
	}
}

type sendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

// Send delivers an email via the Resend API (E01-E05 templates, Section 17 PRD).
// If no API key is configured (e.g. local dev without a Resend account), it logs
// the email instead of failing the caller's request.
func (c *Client) Send(to, subject, html string) error {
	if c.apiKey == "" {
		log.Printf("[email] RESEND_API_KEY not set, skipping send. To=%s Subject=%s", to, subject)
		return nil
	}

	body, err := json.Marshal(sendRequest{
		From:    c.from,
		To:      []string{to},
		Subject: subject,
		HTML:    html,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, resendEndpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("resend API returned status %d", res.StatusCode)
	}
	return nil
}
