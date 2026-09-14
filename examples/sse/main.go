// Command sse streams live businessrules validation events to a browser over
// Server-Sent Events. It demonstrates the fact-producer pattern: the library
// emits events, a listener forwards them onto a go-sse Broadcaster, and the
// /events endpoint fans them out to every connected browser.
//
// Run inside the dev shell (GOEXPERIMENT=jsonv2 is required):
//
//	nix develop --command go run .
//
// Then open http://localhost:8080 and re-run validations by reloading.
package main

import (
	"encoding/json/v2"
	"fmt"
	"log"
	"net/http"
	"time"

	businessrules "github.com/LarsArtmann/go-business-rules"
	"github.com/larsartmann/go-sse"
)

type Order struct {
	Email  string
	Amount float64
	Coupon string
}

func (o Order) rules() []businessrules.Rule {
	return []businessrules.Rule{
		businessrules.Email("email", o.Email, businessrules.SeverityError),
		businessrules.Positive("amount", o.Amount, businessrules.SeverityError),
		businessrules.InRange("amount", o.Amount, 1, 5000, businessrules.SeverityWarning),
		businessrules.OneOf("coupon", o.Coupon, []string{"", "SAVE10", "VIP20"}, businessrules.SeverityInfo),
	}
}

type eventJSON struct {
	Kind      string `json:"kind"`
	Rule      string `json:"rule,omitempty"`
	Severity  string `json:"severity,omitempty"`
	Passed    bool   `json:"passed"`
	Detail    string `json:"detail,omitempty"`
	Valid     *bool  `json:"valid,omitempty"`
	Violation int    `json:"violations,omitempty"`
	Took      string `json:"took"`
}

func main() {
	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", newServer()))
}

func newServer() http.Handler {
	broadcaster := sse.NewBroadcaster[sse.Event]()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, page)
	})

	mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
		stream := sse.NewStream(w, r)
		defer stream.Close()

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
		order := Order{
			Email:  r.FormValue("email"),
			Amount: parseAmount(r.FormValue("amount")),
			Coupon: r.FormValue("coupon"),
		}

		start := time.Now()

		result := businessrules.NewValidator().
			WithListener(func(e businessrules.Event) {
				data := eventJSON{Took: time.Since(start).String()}

				switch ev := e.(type) {
				case businessrules.RuleEvaluated:
					data.Kind = "rule"
					data.Rule = ev.RuleName
					data.Severity = string(ev.Severity)
					data.Passed = ev.Passed()
					if ev.Err != nil {
						data.Detail = ev.Err.Error()
					}
				case businessrules.ValidationCompleted:
					valid := ev.Result.Valid
					data.Kind = "completed"
					data.Valid = &valid
					data.Violation = ev.Result.Count()
				}

				encoded, err := json.Marshal(data)
				if err != nil {
					return
				}

				broadcaster.Broadcast(sse.Event{Event: "validation", Data: string(encoded)})
			}).
			AddRules(order.rules()...).
			Build()

		w.Header().Set("Content-Type", "application/json")

		if err := json.MarshalWrite(w, result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return mux
}

func parseAmount(raw string) float64 {
	var amount float64

	_, _ = fmt.Sscan(raw, &amount)

	return amount
}

const page = `<!doctype html>
<html>
<head><title>businessrules live validation</title>
<style>
 body { font-family: sans-serif; background: #111; color: #eee; margin: 2rem; }
 .event { padding: .4rem .8rem; margin: .3rem 0; border-radius: 6px; background: #222; }
 .pass { border-left: 4px solid #4c8; }
 .fail { border-left: 4px solid #c55; }
 .done { border-left: 4px solid #88f; }
</style>
</head>
<body>
<h1>Live validation events</h1>
<p>Submit a form (or reload) to run validation; every rule check streams in as it happens.</p>
<div id="feed"></div>
<script>
 const feed = document.getElementById('feed');
 const source = new EventSource('/events');
 source.addEventListener('validation', (e) => {
   const data = JSON.parse(e.data);
   const div = document.createElement('div');
   if (data.kind === 'rule') {
     div.className = 'event ' + (data.passed ? 'pass' : 'fail');
     div.textContent = '[' + data.severity + '] ' + data.rule + (data.passed ? ' passed' : ' failed: ' + (data.detail || ''));
   } else {
     div.className = 'event done';
     div.textContent = 'run completed in ' + data.took + ' - ' + (data.valid ? 'valid' : data.violations + ' violation(s)');
   }
   feed.prepend(div);
 });
</script>
</body>
</html>`
