package main

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// The smoke test exercises the full pipeline over real HTTP: a browser-style
// SSE subscription on /events, then a validation run on /validate whose rule
// events must arrive on the stream.
func TestValidationEventsStreamToBrowser(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(newServer())
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventsRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/events", nil)
	if err != nil {
		t.Fatalf("build /events request: %v", err)
	}

	streamed := make(chan string, 16)

	go func() {
		response, err := http.DefaultClient.Do(eventsRequest)
		if err != nil {
			return
		}
		defer response.Body.Close()

		scanner := bufio.NewScanner(response.Body)
		for scanner.Scan() {
			streamed <- scanner.Text()
		}
	}()

	form := url.Values{"email": {"not-an-email"}, "amount": {"-5"}, "coupon": {"SAVE10"}}
	validateRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, server.URL+"/validate", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("build /validate request: %v", err)
	}

	validateRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validateResponse, err := http.DefaultClient.Do(validateRequest)
	if err != nil {
		t.Fatalf("POST /validate: %v", err)
	}
	defer validateResponse.Body.Close()

	if validateResponse.StatusCode != http.StatusOK {
		t.Fatalf("POST /validate status = %d, want 200", validateResponse.StatusCode)
	}

	var sawRuleEvent, sawCompletedEvent bool

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("stream closed early; saw rule event: %v, completed event: %v", sawRuleEvent, sawCompletedEvent)
		case line := <-streamed:
			switch {
			case strings.Contains(line, `"kind":"rule"`):
				sawRuleEvent = true
			case strings.Contains(line, `"kind":"completed"`):
				sawCompletedEvent = true
			}

			if sawRuleEvent && sawCompletedEvent {
				return
			}
		}
	}
}
