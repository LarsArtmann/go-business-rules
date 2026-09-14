module github.com/LarsArtmann/go-business-rules/listeners/otel

go 1.26.7

require (
	github.com/LarsArtmann/go-business-rules/v2 v2.1.0
	github.com/onsi/ginkgo/v2 v2.32.2
	github.com/onsi/gomega v1.43.0
	go.opentelemetry.io/otel v1.46.0
	go.opentelemetry.io/otel/metric v1.46.0
	go.opentelemetry.io/otel/sdk v1.46.0
	go.opentelemetry.io/otel/sdk/metric v1.46.0
	go.opentelemetry.io/otel/trace v1.46.0
)

require (
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20260709232956-b9395ee17fa0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/larsartmann/go-error-family v0.10.0 // indirect
	github.com/larsartmann/go-finding v1.10.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/mod v0.41.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	golang.org/x/tools v0.49.0 // indirect
)

replace github.com/LarsArtmann/go-business-rules/v2 => ../..
