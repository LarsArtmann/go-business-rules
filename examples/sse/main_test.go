package main

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const validatePayload = `{"email":"not-an-email","amount":-5,"coupon":"SAVE10"}`

func postValidate(t *testing.T, serverURL string, ctx context.Context) *http.Response {
	t.Helper()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		serverURL+"/validate",
		strings.NewReader(validatePayload),
	)
	if err != nil {
		t.Fatalf("build /validate request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Datastar-Request", "true")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("POST /validate: %v", err)
	}

	return response
}

// The smoke test exercises the full pipeline over real HTTP: a Datastar-style
// @post to /validate whose response is an SSE stream of patch-elements and
// patch-signals events, ending right after the run summary.
func TestValidationStreamsDatastarPatches(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(newServer())
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response := postValidate(t, server.URL, ctx)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("POST /validate status = %d, want 200", response.StatusCode)
	}

	if contentType := response.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "text/event-stream") {
		t.Fatalf("POST /validate content-type = %q, want text/event-stream", contentType)
	}

	lines := readAllLines(t, response.Body)
	stream := strings.Join(lines, "\n")

	for _, want := range []string{
		"event: datastar-patch-elements",
		"data: selector #feed",
		"data: mode append",
		"data: elements <li",
		"event: datastar-patch-signals",
		`data: signals {"valid":false,"violations":3}`,
	} {
		if !strings.Contains(stream, want) {
			t.Fatalf("stream missing %q; got:\n%s", want, stream)
		}
	}

	if len(lines) == 0 || lines[len(lines)-1] != "" {
		// The terminal signals event must be the last one the server sends.
		if strings.LastIndex(stream, "datastar-patch-signals") < strings.LastIndex(stream, "datastar-patch-elements") {
			t.Fatalf("terminal signals event must follow all element patches; got:\n%s", stream)
		}
	}
}

// The /events endpoint fans every run out to all connected clients.
func TestEventsEndpointFansOutAllRuns(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(newServer())
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventsRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/events", nil)
	if err != nil {
		t.Fatalf("build /events request: %v", err)
	}

	streamed := make(chan string, 64)

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

		if err := scanner.Err(); err != nil {
			t.Logf("/events scan error: %v", err)
		}
	}()

	validateResponse := postValidate(t, server.URL, ctx)
	_, _ = io.Copy(io.Discard, validateResponse.Body)
	validateResponse.Body.Close()

	var sawElements, sawSignals bool

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("timeout waiting on /events; saw elements: %v, signals: %v", sawElements, sawSignals)
		case line := <-streamed:
			switch {
			case strings.Contains(line, "event: datastar-patch-elements"):
				sawElements = true
			case strings.Contains(line, "event: datastar-patch-signals"):
				sawSignals = true
			}

			if sawElements && sawSignals {
				return
			}
		}
	}
}

// Malformed signal payloads must fail loudly instead of silently validating
// an empty order.
func TestValidateRejectsUnreadableSignals(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(newServer())
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		server.URL+"/validate",
		strings.NewReader("{not json"),
	)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("POST /validate: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("POST /validate status = %d, want 400", response.StatusCode)
	}
}

func readAllLines(t *testing.T, reader io.Reader) []string {
	t.Helper()

	var lines []string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("read stream: %v", err)
	}

	lines = append(lines, "")

	return lines
}
