// Command sse streams live businessrules validation events to a browser as
// Datastar patches. It demonstrates the fact-producer pattern twice over:
// the library emits events, a listener turns them into Datastar
// patch-elements / patch-signals values, and those values are broadcast
// through one go-sse Broadcaster. POST /validate streams the patches for its
// own run back to the browser as the response; GET /events fans every run
// out to any other connected client (e.g. curl or a dashboard).
//
// Run inside the dev shell (GOEXPERIMENT=jsonv2 is required):
//
//	nix develop --command go run .
//
// Then open http://localhost:8080, fill the form, and watch rule checks
// stream in as they happen.
package main

import (
	"encoding/json/v2"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/larsartmann/go-sse"
	"github.com/starfederation/datastar-go/datastar"

	businessrules "github.com/LarsArtmann/go-business-rules/v2"
)

type Order struct {
	Email  string  `json:"email"`
	Amount float64 `json:"amount"`
	Coupon string  `json:"coupon"`
}

func (o Order) rules() []businessrules.Rule {
	return []businessrules.Rule{
		businessrules.Email("email", o.Email, businessrules.SeverityError),
		businessrules.Positive("amount", o.Amount, businessrules.SeverityError),
		businessrules.InRange("amount", o.Amount, 1, 5000, businessrules.SeverityWarning),
		businessrules.OneOf("coupon", o.Coupon, []string{"", "SAVE10", "VIP20"}, businessrules.SeverityInfo),
	}
}

// summarySignals is patched into the browser when a validation run finishes.
type summarySignals struct {
	Valid      bool `json:"valid"`
	Violations int  `json:"violations"`
}

// datastarRetry mirrors datastar.DefaultSseRetryDuration, in milliseconds,
// the unit go-sse writes into the retry field.
const datastarRetry = 1000

// elementsPatch renders a Datastar patch-elements SSE event as a plain
// value, using the exact data-line protocol from the Datastar reference.
func elementsPatch(eventID sse.EventID, selector string, mode datastar.ElementPatchMode, fragment string) sse.Event {
	rows := make([]string, 0, 3)

	rows = append(rows, datastar.SelectorDatalineLiteral+selector)

	if mode != datastar.ElementPatchModeOuter {
		rows = append(rows, datastar.ModeDatalineLiteral+string(mode))
	}

	for line := range strings.SplitSeq(fragment, "\n") {
		rows = append(rows, datastar.ElementsDatalineLiteral+line)
	}

	return sse.Event{
		Event: string(datastar.EventTypePatchElements),
		ID:    eventID,
		Retry: datastarRetry,
		Data:  strings.Join(rows, "\n"),
	}
}

// signalsPatch renders a Datastar patch-signals SSE event as a plain value.
func signalsPatch(eventID sse.EventID, payload []byte) sse.Event {
	return sse.Event{
		Event: string(datastar.EventTypePatchSignals),
		ID:    eventID,
		Retry: datastarRetry,
		Data:  datastar.SignalsDatalineLiteral + string(payload),
	}
}

// ruleRow renders one feed entry for a single rule evaluation.
func ruleRow(ev businessrules.RuleEvaluated) string {
	class, detail := "pass", ""

	if !ev.Passed() {
		class = "fail"
		detail = ": " + html.EscapeString(ev.Err.Error())
	}

	return fmt.Sprintf(
		`<li class=%q>[%s] %s%s</li>`,
		class,
		html.EscapeString(string(ev.Severity)),
		html.EscapeString(ev.RuleName),
		detail,
	)
}

// runValidation evaluates the order and turns every validation event into a
// Datastar patch tagged with the run ID.
func runValidation(order Order, eventID sse.EventID, broadcast func(sse.Event)) {
	businessrules.NewValidator().
		WithListener(func(e businessrules.Event) {
			switch ev := e.(type) {
			case businessrules.RuleEvaluated:
				broadcast(elementsPatch(eventID, "#feed", datastar.ElementPatchModeAppend, ruleRow(ev)))
			case businessrules.ValidationCompleted:
				signals, err := json.Marshal(summarySignals{
					Valid:      ev.Result.Valid,
					Violations: ev.Result.Count(),
				})
				if err == nil {
					broadcast(signalsPatch(eventID, signals))
				}
			}
		}).
		AddRules(order.rules()...).
		Build()
}

func main() {
	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", newServer()))
}

func newServer() http.Handler {
	broadcaster := sse.NewBroadcaster[sse.Event](sse.WithBufferSize[sse.Event](64))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if _, err := fmt.Fprint(w, page); err != nil {
			log.Printf("write index page: %v", err)
		}
	})

	mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
		stream := sse.NewStream(w, r)
		defer func() { _ = stream.Close() }()

		events := broadcaster.Subscribe()
		defer broadcaster.Unsubscribe(events)

		for {
			select {
			case <-stream.Context().Done():
				return
			case evt, ok := <-events:
				if !ok || stream.Send(evt) != nil {
					return
				}
			}
		}
	})

	mux.HandleFunc("POST /validate", func(w http.ResponseWriter, r *http.Request) {
		var order Order

		if err := datastar.ReadSignals(r, &order); err != nil {
			http.Error(w, "could not read Datastar signals: "+err.Error(), http.StatusBadRequest)

			return
		}

		eventID, err := sse.ParseEventID(fmt.Sprintf("run-%d", time.Now().UnixNano()))
		if err != nil {
			http.Error(w, "invalid run id: "+err.Error(), http.StatusInternalServerError)

			return
		}

		stream := sse.NewStream(w, r)
		defer func() { _ = stream.Close() }()

		events := broadcaster.Subscribe()
		defer broadcaster.Unsubscribe(events)

		finished := make(chan struct{})

		go func() {
			defer close(finished)

			runValidation(order, eventID, broadcaster.Broadcast)
		}()

		forward := func(evt sse.Event) bool {
			if evt.ID != eventID {
				return false
			}

			if stream.Send(evt) != nil {
				return true
			}

			return evt.Event == string(datastar.EventTypePatchSignals)
		}

		for {
			select {
			case <-r.Context().Done():
				return
			case evt, ok := <-events:
				if !ok || forward(evt) {
					return
				}
			case <-finished:
				for {
					select {
					case evt, ok := <-events:
						if !ok || forward(evt) {
							return
						}
					default:
						return
					}
				}
			}
		}
	})

	return mux
}

// datastarBundleVersion must match the version of the Go SDK above, so the
// wire protocol and the client runtime stay in lockstep.
const datastarBundleVersion = "v1.2.2"

const page = `<!doctype html>
<html>
<head><title>businessrules live validation</title>
<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@` + datastarBundleVersion + `/bundles/datastar.js"></script>
<style>
 body { font-family: sans-serif; background: #111; color: #eee; margin: 2rem; }
 form { display: flex; gap: 1rem; align-items: end; margin: 1rem 0; }
 label { display: flex; flex-direction: column; font-size: .8rem; }
 input { background: #222; color: #eee; border: 1px solid #444; border-radius: 6px; padding: .4rem; }
 button { background: #4c8; border: 0; border-radius: 6px; padding: .5rem 1rem; font-weight: bold; }
 .event { padding: .4rem .8rem; margin: .3rem 0; border-radius: 6px; background: #222; list-style: none; }
 .pass { border-left: 4px solid #4c8; }
 .fail { border-left: 4px solid #c55; }
</style>
</head>
<body data-signals="{valid: '', violations: 0, amount: 0}">
<h1>Live validation events</h1>
<form data-on:submit="@post('/validate')">
 <label>email <input data-bind:email /></label>
 <label>amount <input type="number" step="0.01" data-bind:amount /></label>
 <label>coupon <input data-bind:coupon /></label>
 <button>Validate</button>
</form>
<p data-text="$valid === '' ? 'Submit to run validation' : ($valid ? 'valid' : $violations + ' violation(s)')"></p>
<ol id="feed"></ol>
</body>
</html>`
