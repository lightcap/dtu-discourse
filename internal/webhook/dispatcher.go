package webhook

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// GamificationPayload matches Eve's Zod schema for gamification webhooks.
type GamificationPayload struct {
	DiscourseUserID             int         `json:"discourseUserId"`
	Action                      string      `json:"action"`
	DiscourseResourceID         interface{} `json:"discourseResourceId"`
	CounterpartyDiscourseUserID *int        `json:"counterpartyDiscourseUserId,omitempty"`
	OccurredAt                  string      `json:"occurredAt,omitempty"`
}

// Dispatcher sends webhook payloads to a configured URL.
type Dispatcher struct {
	URL    string
	Secret string
	client *http.Client
}

// New creates a Dispatcher. If url is empty, Dispatch is a no-op.
func New(url, secret string) *Dispatcher {
	return &Dispatcher{
		URL:    url,
		Secret: secret,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// Dispatch sends the payload in a fire-and-forget goroutine.
func (d *Dispatcher) Dispatch(payload GamificationPayload) {
	if d == nil || d.URL == "" {
		return
	}
	if payload.OccurredAt == "" {
		payload.OccurredAt = time.Now().UTC().Format(time.RFC3339)
	}
	go func() {
		body, err := json.Marshal(payload)
		if err != nil {
			log.Printf("[webhook] marshal error: %v", err)
			return
		}
		req, err := http.NewRequest("POST", d.URL, bytes.NewReader(body))
		if err != nil {
			log.Printf("[webhook] request error: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		if d.Secret != "" {
			req.Header.Set("x-webhook-secret", d.Secret)
		}
		resp, err := d.client.Do(req)
		if err != nil {
			log.Printf("[webhook] dispatch error: %v", err)
			return
		}
		resp.Body.Close()
		log.Printf("[webhook] %s dispatched to %s — %d", payload.Action, d.URL, resp.StatusCode)
	}()
}
